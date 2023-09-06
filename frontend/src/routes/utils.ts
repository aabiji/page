export async function call_api(url: string, payload: object = {}): Promise<any> {
    const response = await fetch(url, {
        method: "GET", // change to POST
        mode: "cors",
        //body: JSON.stringify(payload),
    });
    return response.text(); // change to json
}

export async function download_file(url: string): Promise<string> {
    const response = await fetch(url, {
        method: "GET",
        mode: "cors",
    });
    return response.text();
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
