extern crate reqwest;
extern crate regex;

use std::io::Read;
use regex::Regex;

fn fetch(url: &str) -> Result<Option<String>,String> {
    let mut resp = try!(reqwest::get(url).map_err(|e| e.to_string()));
    let mut content = String::new();
    match resp.read_to_string(&mut content).map_err(|e| e.to_string()) {
        Ok(_) => { return Ok(Some(content)) },
        _ => { return Ok(None) }
    }
}

fn extract_links(content: &str) -> Vec<String> {
    let re = Regex::new("(href|src)=\"(\\S+)\"").unwrap();
    return re.captures_iter(content).map(|cap| String::from(&cap[2])).collect();
}

#[derive(Debug)]
enum LinkType {
    Relative,
    Absolute,
    External,
    Ignored,
}

fn classify_link(link: &str) -> LinkType {
    let absolute = Regex::new("^/").unwrap();
    if absolute.is_match(link) {
        return LinkType::Absolute;
    }
    let external = Regex::new("^https?://").unwrap();
    if external.is_match(link) {
        return LinkType::External;
    }
    let ignored = Regex::new("^(mailto)").unwrap();
    if ignored.is_match(link) {
        return LinkType::Ignored;
    }

    return LinkType::Relative;
}

fn format_relative_link(link: &str, current_url: &str) -> String {
    return format!("{}{}", current_url, link);
}

fn format_absolute_link(link: &str, current_url: &str) -> String {
    let domain_matcher = Regex::new("(https?://[^/]+)").unwrap();
    match domain_matcher.find(current_url) {
        Some(domain) => {
            return format!("{}{}", domain.as_str(), link);
        },
        None => { println!("failed"); return String::from(link) }
    }
}

fn content_links(content: &str, url: &str) -> Vec<String> {
    let mut links = vec![];
    for link in extract_links(&content) {
        match classify_link(link.as_str()) {
            LinkType::Relative => {
                links.push(format_relative_link(link.as_str(), url));
            },
            LinkType::Absolute => {
                links.push(format_absolute_link(link.as_str(), url));
            },
            LinkType::External => {
                links.push(String::from(link));
            },
            LinkType::Ignored => { println!("SKIP\t {}", link); }
        }
    }
    return links;
}

fn process_page(url: &str) -> Option<Vec<String>> {
    match fetch(url) {
        Ok(c) => {
            println!("OK\t {}", url);
            if c.is_some() {
                return Some(content_links(c.unwrap().as_str(), url));
            } else {
                return None;
            }
        },
        Err(_) => {
            println!("ERROR\t {}", url);
            return None;
        }
    };
}

fn main() {
    let mut links = vec![String::from("https://charlieegan3.com")];
    process_page("https://charlieegan3.com/");
}
