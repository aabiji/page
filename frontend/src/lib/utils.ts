import { goto } from "$app/navigation";

export const serverError   = "Server error";
export const backendOrigin = "http://localhost:8080";

// Key names used to store user data in localStorage
export const BooksKey    = "User:Books";
export const SettingsKey = "User:Settings";
export const BookKey     = (id: number) => `Book:${id}`;
export const UserBookKey = (id: number) => `Userbook:${id}`;

// JSON structure for book info returned by backend API at endpoint /book/get/{id}
interface Section { Name: string, Path: string };
interface Info {
    Author: string,     Title: string,    Contributor: string,
    Coverage: string,   Date: string,     Description: string,
    Identifier: string, Language: string, Publisher: string,
    Relation: string,   Rights: string,   Source: string,
    Subjects: string[],
};
export class Book {
    CoverImagePath: string = "";
    TableOfContents: Section[] = [{Name: "", Path: ""}];
    Info: Info = {Author: "", Title: "", Contributor: "", Coverage: "",
                  Date: "", Description: "", Identifier: "", Language: "",
                  Publisher: "", Relation: "", Rights: "", Source: "", Subjects: [""]};
}

export async function callApi(url: string, method: string, json: object = {}, isFile: boolean=false): Promise<any> {
    let data = isFile ? json : JSON.stringify(json);
    let payload = {
        method: method,
        credentials: "include",
        body: method == "POST" ? data : null,
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

// Redirect to auth page if user has not authenticated
export function redirectIfNotAuth() {
    if (document.cookie == "") {
        goto("/auth");
    }
}

export function staticFileUrl(file: string): string {
    file = file.replace(window.location.origin+"/", "");
    return `${backendOrigin}/static/${file}`;
}

export function coverImagePath(file: string): string {
    return file == "" ? "default-cover-image.png" : staticFileUrl(file);
}

export function cacheGet(key: string): any {
    let obj = localStorage.getItem(key);
    let type = key == BooksKey ? [] : {};
    return obj == null ? type : JSON.parse(obj);
}

export function cacheBook(id: number, info: object) {
    localStorage.setItem(BookKey(id), JSON.stringify(info));
    let bookIds = cacheGet(BooksKey) as number[];
    bookIds.push(id);
    localStorage.setItem(BooksKey, JSON.stringify(bookIds));
}

export function removeCookie(name: string) {
    document.cookie = `${name}=; Max-Age=-99999999;`;
}