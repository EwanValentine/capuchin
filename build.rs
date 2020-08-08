fn main() {
    protoc_rust_grpc::Codegen::new()
        .out_dir("src/node/leader")
        .input("proto/leader.proto")
        .rust_protobuf(true)
        .run()
        .expect("protoc-rust-grpc");
}
