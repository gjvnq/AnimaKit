#[macro_use]
extern crate clap;
use clap::App;

extern crate sdl2; 
use sdl2::pixels::Color;
use sdl2::event::Event;
use sdl2::keyboard::Keycode;
use sdl2::render::Canvas;
use sdl2::video::Window;
use std::time::Duration;

fn draw_time_controls(canvas: &mut Canvas<Window>) {
	// Get canvas size
	let full_rect = canvas.output_size();
	let (x, y) = full_rect.unwrap();
	println!("{} {}", x, y);
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