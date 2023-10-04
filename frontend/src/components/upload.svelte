<script lang="ts">
    import * as utils from "$lib/utils";
 
    function addBook(id: number) {
        let url = `${utils.backendOrigin}/book/get/${id}`
        utils.callApi(url, "GET").then((info: utils.Book) => {
            info.CoverImagePath = utils.coverImagePath(info.CoverImagePath);
            utils.cacheBook(id, info);
        });
    }

    function uploadFile(file: File): Promise<void> {
        return new Promise((resolve, reject) => {
            const formData = new FormData();
            formData.append("file", file);
            let url = `${utils.backendOrigin}/user/book/upload`;
            utils.callApi(url, "POST", formData, true).then((response) => {
                if (utils.serverError in response) {
                    console.log(response);
                    reject();
                    return;    
                }
                addBook(response.BookId);
                resolve();
            });
        });
    }

    let fileInput: HTMLElement;
    function uploadSelectedFile(event: any) {
        const file = event.target.files[0];
        uploadFile(file).then(() => {
            window.location.pathname = "/";
        });
    }

    let uploadPromises: Promise<void>[] = [];
    const disableDrag = (event: DragEvent) => event.preventDefault();
    function uploadDroppedFile(event: DragEvent) {
        event.preventDefault();
        let files = event.dataTransfer?.files;
        if (files == undefined) return;
        for (let i = 0; i < files.length; i++) {
            let file = files[i];
            let extension = file.name.split(".").at(-1)
            if (extension != "epub") continue;
            uploadPromises.push(uploadFile(file));
        }
        Promise.all(uploadPromises).then(() => {
            uploadPromises = [];
            window.location.pathname = "/";
        });
    }

    const openFileDialog = () => fileInput.click();
</script>

<dialog open>
    <h1> Upload a book </h1>
    <input on:change={uploadSelectedFile} bind:this={fileInput} accept=".epub" type="file" style="display: none;">
    <div tabindex=-1 aria-label="Click, Drop, and Drag Over Element" role="button" class="drop" 
         on:keypress={openFileDialog} on:click={openFileDialog} on:drop={uploadDroppedFile} on:dragover={disableDrag}>
        <p> Click me! </p>
        <p> OR </p>
        <p> Drag your epub files here. </p>
    </div>
</dialog>

<style>
    dialog {
        width: 450px;
        height: 300px;
        text-align: center;
        color: white;
        border: none;
        position: absolute;
        margin-top: 50px;
        z-index: 2;
        background-color: var(--background-accent);
        box-shadow: rgba(0, 0, 0, 0.56) 0px 22px 70px 4px;
    }

    .drop {
        height: 200px;
        border: 3px dashed #6b6b6b;
        background-color: #2b2c2c;
    }
</style>