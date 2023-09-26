import { goto } from "$app/navigation";

export const serverError = "Server error";

export async function callApi(url: string, method: string, json: object): Promise<any> {
    let payload = {
        method: method,
        credentials: "include",
        body: method == "POST" ? JSON.stringify(json) : null,
    };
    const response = await fetch(url, payload as RequestInit);
    return response.json();
}

export async function downloadFile(url: string): Promise<string> {
    const response = await fetch(url, {
        method: "GET",
        mode: "cors",
    });
    return response.text();
}

export async function hashSHA256(data: string): Promise<string> {
    let encoded = new TextEncoder().encode(data);
    let buffer = await window.crypto.subtle.digest("SHA-256", encoded);
    let hash = Array.from(new Uint8Array(buffer));
    return hash.map(byte => byte.toString(16).padStart(2, "0")).join("");
}

export function staticFileUrl(file: string): string {
    file = file.replace(window.location.origin+"/", "");
    return `http://localhost:8080/static/${file}`;
}

export function coverImagePath(file: string): string {
    return file == "" ? "default-cover-image.png" : staticFileUrl(file);
}

// Redirect to auth page if user has not authenticated
export function redirectIfNotAuth() {
    if (document.cookie == "") {
        goto("/auth");
    }
}
