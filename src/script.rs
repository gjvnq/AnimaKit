use v8::{self, value};

fn load_script_from_file(path: &str) Result<> {
	let isolate = v8::Isolate::new();
	let context = v8::Context::new(&isolate);
}