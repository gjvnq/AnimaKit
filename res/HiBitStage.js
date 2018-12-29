function HiBitStage(width, height) {
  var id = ffi_HiBitStage_new(width, height);
  return {
    get id() { return id },
    get width() { return width },
    get height() { return height },
    get fps() { return fps },
    get bg() { return ffi_HiBitStage_get_bg(id) }
    set bg(new_bg) { return ffi_HiBitStage_set_bg(id, new_bg) },
    place: function (obj) {}
  };
}
