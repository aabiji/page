export const serverError = "Server error";

export async function callApi(url: string, method: string, json: object): Promise<any> {
    let payload = {
        method: method,
        body: method == "POST" ? JSON.stringify(json) : null,
        credentials: "include",
    };
    const response = await fetch(url, payload);
    return response.json();
}

export async function downloadFile(url: string): Promise<string> {
    const response = await fetch(url, {
        method: "GET",
        mode: "cors",
    });
    return response.text();
}

export function staticFileUrl(file: string): string {
    file = file.replace(window.location.origin+"/", "");
    return `http://localhost:8080/static/${file}`;
}

export function coverImagePath(file: string): string {
    return file == "" ? "default-cover-image.png" : staticFileUrl(file);
}
