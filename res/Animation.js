var Animation = {
  _stage: null,
  
  get stage() { return this._stage }
  set stage(new_stage) { this._stage = new_stage; ffi_Animation_set_stage(new_stage._id) }

  get width() { return ffi_Animation_get_width() }
  set width(new_width) { return ffi_Animation_set_width(new_width) }
  get height() { return ffi_Animation_get_height() }
  set height(new_height) { return ffi_Animation_set_height(new_height) }
  get length() { return ffi_Animation_get_length() }
  set length(new_length) { return ffi_Animation_set_length(new_length) }
  get fps() { return ffi_Animation_get_fps() }
  set fps(new_fps) { return ffi_Animation_set_fps(new_fps) }  
}