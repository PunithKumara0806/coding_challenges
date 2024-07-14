use std::{
    env,
    fs::File,
    io::{self, BufRead},
};

fn main() {
    let file_1 = env::args().nth(1).expect("No file name found");
    let file_2 = env::args().nth(2).expect("No file name found");
    let fd_1 = File::open(&file_1).expect(&format!("Can't open file :{}", file_1));
    let fd_2 = File::open(&file_2).expect(&format!("Can't open file :{}", file_2));
    let stream_1 = io::BufReader::new(fd_1);
    let stream_2 = io::BufReader::new(fd_2);
    let mut line = String::new();
}
