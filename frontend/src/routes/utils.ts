export async function call_api(url: string, text: boolean = true) {
    const response = await fetch(url, {
        method: "GET",
        mode: "cors",
    });
    return text ? await response.text() : response.json();
}

export function get_content_type(file_url: string): DOMParserSupportedType {
    // all content types:
    // ["application/xhtml+xml","application/xml","image/svg+xml","text/html","text/xml"]
    let extention = file_url.split(".")[1];
    return extention == "xhtml" ? "application/xhtml+xml" : "text/html";
}

export function static_file_url(file: string): string {
    let netpath = file.replace("BOOKS/", "");
    return `http://localhost:8080/static/${netpath}`;
}

export function book_info_url(book_name: string): string {
    return `http://localhost:8080/${book_name}`
}
