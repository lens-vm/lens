import { JSON, JSONEncoder } from "assemblyscript-json";

@external("lens", "next")
export declare function next(): usize

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

class StreamOption {
    json: String
    endOfStream: bool

    constructor(json: String, endOfStream: bool) {
        this.json = json;
        this.endOfStream = endOfStream;
    }
}

const JSON_TYPE_ID: i8 = 1;
const EOS_TYPE_ID: i8 = 127;

export function alloc(size: usize): usize {
    let buf = new ArrayBuffer(i32(size));
    let ptr = changetype<usize>(buf);
    store<ArrayBuffer>(ptr, buf);
    return ptr;
}

export function transform(): usize {
    let ptr = next();
    let streamOption = fromTransportVec(ptr)
    if (streamOption.endOfStream) {
        return toTransportVec(EOS_TYPE_ID, "");
    }

    let inputJsonObj = <JSON.Obj>(JSON.parse(streamOption.json));

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

function fromTransportVec(ptr: usize): StreamOption {
    let type = load<i8>(ptr)
    if (type == EOS_TYPE_ID) {
        return new StreamOption("", true)
    }
    let len = load<u32>(ptr+1)
    return new StreamOption(String.UTF8.decodeUnsafe(ptr+1+4, len, false), false)
}

function toTransportVec(type_id: i8, message: string): usize {
    let len = String.UTF8.byteLength(message, false);

    let buf = new Uint8Array(len+1+4);
    let ptr = changetype<usize>(buf);
    store<i8>(ptr, type_id);
    if (type_id == EOS_TYPE_ID) {
        return ptr
    }
    store<u32>(ptr+1, len);

    String.UTF8.encodeUnsafe(changetype<usize>(message), len, ptr+1+4, false)

    return ptr
}
