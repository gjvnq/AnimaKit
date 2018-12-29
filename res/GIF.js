function GIF(filename) {
  var id = ffi_GIF_new(filename);
  return {
    get id() { return id },
    get frames() { return ffi_GIF_get_frames(id) },
    get keyframes() { return ffi_GIF_get_keyframes(id) }
    set keyframes(new_keys) { return ffi_GIF_set_keyframes(id, new_keys) }
  };
}
