function Image(filename) {
  var id = ffi_Image_new(filename);
  return {
    get id() { return id },
    get keyframes() { return ffi_Image_get_keyframes(id) }
    set keyframes(new_keys) { return ffi_Image_set_keyframes(id, new_keys) }
  };
}
