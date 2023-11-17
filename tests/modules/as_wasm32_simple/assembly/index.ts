// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

import { JSON, JSONEncoder } from "assemblyscript-json"; 

// AssemblyScript demands that an abort func exists, so we define our own here and ask the compiler to use it
// instead of an import (see asconfig.json "use" flag), and https://www.assemblyscript.org/concepts.html#special-imports
function abort(
    message: string | null,
    fileName: string | null,
    lineNumber: u32,
    columnNumber: u32
  ): void {
    // This just interupts the wasm execution, returning a default value. We should
    // do this better when we build a serious AssemblyScript SDK
    unreachable()
}

const JSON_TYPE_ID: i8 = 1;

export function alloc(size: usize): usize {
    return heap.alloc(size);
}

export function free(ptr: usize, size: usize): void {
    heap.free(ptr)
}

export function transform(ptr: usize): usize {
    let inputStr = fromTransportVec(ptr)
    let inputJsonObj = <JSON.Obj>(JSON.parse(inputStr));

    let inputName = inputJsonObj.getString("FullName");
    let inputAge = inputJsonObj.getInteger("Age");

    let encoder = new JSONEncoder();
    encoder.pushObject(null)
    if (inputName != null) {
        encoder.setString("FullName", inputName.valueOf())
    }
    if (inputAge != null) {
        encoder.setInteger("Age", inputAge.valueOf() + 10)
    }
    encoder.popObject()

    return toTransportVec(JSON_TYPE_ID, encoder.toString());
}

function fromTransportVec(ptr: usize): string {
    let type = load<i8>(ptr)
    let len = load<u32>(ptr+1)
    return String.UTF8.decodeUnsafe(ptr+1+4, len, false)
}

function toTransportVec(type_id: i8, message: string): usize {
    let len = String.UTF8.byteLength(message, false);

    let buf = new Uint8Array(len+1+4);
    let ptr = changetype<usize>(buf);
    store<i8>(ptr, len);
    store<u32>(ptr+1, len);

    String.UTF8.encodeUnsafe(changetype<usize>(message), len, ptr+1+4, false)

    return ptr
}
