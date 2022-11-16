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

export function alloc(size: usize): usize {
    let buf = new ArrayBuffer(i32(size));
    let ptr = changetype<usize>(buf);
    store<ArrayBuffer>(ptr, buf);
    return ptr;
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

    return toTransportVec(encoder.toString());
}

function fromTransportVec(ptr: usize): string {
    let len = load<u32>(ptr)
    return String.UTF8.decodeUnsafe(ptr+4, len, false)
}

function toTransportVec(message: string): usize {
    let len = String.UTF8.byteLength(message, false);

    let buf = new Uint8Array(len+4);
    let ptr = changetype<usize>(buf);
    store<u32>(ptr, len);

    String.UTF8.encodeUnsafe(changetype<usize>(message), len, ptr+4, false)

    return ptr
}
