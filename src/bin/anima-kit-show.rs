#[macro_use]
extern crate clap;
use clap::App;

extern crate sdl2; 
extern crate sdl2_unifont;
use sdl2::pixels::Color;
use sdl2::event::Event;
use sdl2::rect::Rect;
use sdl2::keyboard::Keycode;
use sdl2::render::Canvas;
use sdl2::video::Window;
use std::time::Duration;
use sdl2_unifont::renderer::SurfaceRenderer as FontSurfaceRenderer;


pub const TRANSPARENT: Color = Color{r: 0, g: 0, b: 0, a: 0};
pub const WHITE: Color = Color{r: 255, g: 255, b: 255, a: 0};
pub const BLUE: Color = Color{r: 0, g: 122, b: 255, a: 0};
pub const BROWN: Color = Color{r: 162, g: 132, b: 94, a: 0};
pub const GREY: Color = Color{r: 142, g: 142, b: 147, a: 0};
pub const GREEN: Color = Color{r: 40, g: 205, b: 65, a: 0};
pub const ORANGE: Color = Color{r: 255, g: 149, b: 0, a: 0};
pub const PINK: Color = Color{r: 255, g: 45, b: 85, a: 0};
pub const PURPLE: Color = Color{r: 175, g: 82, b: 222, a: 0};
pub const RED: Color = Color{r: 255, g: 59, b: 48, a: 0};
pub const YELLOW: Color = Color{r: 255, g: 204, b: 0, a: 0};

#[allow(dead_code)]
static mut TIME_POS : f64 = 0.0;


fn draw_time_controls(canvas: &mut Canvas<Window>) {
	// Get canvas size
	let full_rect = canvas.output_size();
	let (window_w, window_h) = full_rect.unwrap();

	// Compute the controls height
	let controls_width = window_w;
	let controls_height = 64;

	// Draw base
	canvas.set_draw_color(GREY);
	canvas.fill_rect(Rect::new(0, (window_h - controls_height) as i32, controls_width, controls_height)).unwrap();

	// Write instructions
	let texture_creator = canvas.texture_creator();
	let txt_render = FontSurfaceRenderer::new(RED, BLUE);
	let _keys_txt = "[SPACE] play/pause and [←] [→] to move frame by frame";

	// txt_render
	// 	.draw(keys_txt)
	// 	.unwrap()
	// 	.blit(None, &mut screen, Rect::new(2, 2, 0, 0))
	// 	.unwrap();
	let keys_surface = txt_render.draw("hi").unwrap();
	let keys_texture = texture_creator.create_texture_from_surface(keys_surface).unwrap();
	canvas.copy(&keys_texture, None, None).unwrap();
}

fn start_ui(filename: Option<&str>) {
	let sdl_context = sdl2::init().unwrap();
    let video_subsystem = sdl_context.video().unwrap();
 
 	let window_title = match filename {
 		None => "AnimaKit".to_owned(),
 		Some(f) => format!("AnimaKit - {}", f)
 	};
    let window = video_subsystem.window(window_title.as_str(), 800, 600)
        .position_centered()
        .resizable()
        .build()
        .unwrap();
 
    let mut canvas = window.into_canvas().build().unwrap();
 
    let mut event_pump = sdl_context.event_pump().unwrap();
    'main_loop: loop {
        // Clear the screen to black
        canvas.set_draw_color(Color::RGB(0, 0, 0));
        canvas.clear();
    	// Process events
        for event in event_pump.poll_iter() {
            match event {
                Event::Quit {..} |
                Event::KeyDown { keycode: Some(Keycode::Escape), .. } => {
                    break 'main_loop
                },
                _ => {}
            }
        }
        // Process stuff
        draw_time_controls(&mut canvas);

        // Show screen and wait
        canvas.present();
        ::std::thread::sleep(Duration::new(0, 1_000_000_000u32 / 60));
    }
}

fn main() {
    // The YAML file is found relative to the current file, similar to how modules are found
    let yaml = load_yaml!("show-cli.yml");
    let matches = App::from_yaml(yaml).get_matches();

    // Same as previous examples...
    let input = matches.value_of("INPUT").unwrap();
    println!("Using input file: {}", input);

    start_ui(Some(input));
}