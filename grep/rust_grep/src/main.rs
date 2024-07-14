use std::{fs::File, io::Stdin, os::fd::AsFd};

use clap::Parser;

#[derive(Parser)]
#[command(version,about = "rust_grep, a simple grep written in rust",long_about = None)]
struct Cli {
    pattern: String,
    /// File path, if not given takes stdin by default
    file: Option<String>,
}

fn main() {
    let cli = Cli::parse();
    let path = match cli.file {
        Some(v) => v,
        None => "".to_string(),
    };
}
