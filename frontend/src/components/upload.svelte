<script lang="ts">
    import * as utils from "$lib/utils";
 
    function addBook(id: number) {
        let url = `${utils.backendOrigin}/book/get/${id}`
        utils.callApi(url, "GET").then((info: utils.Book) => {
            info.CoverImagePath = utils.coverImagePath(info.CoverImagePath);
            utils.cacheBook(id, info);
        });
    }

    function uploadFile(file: File) {
        const formData = new FormData();
        formData.append("file", file);
        let url = `${utils.backendOrigin}/user/book/upload`;
        utils.callApi(url, "POST", formData, true).then((response) => {
            if (utils.serverError in response) {
                console.log(response);
                return;    
            }
            addBook(response.BookId);
        });
    }

    let fileInput: HTMLElement;
    function uploadSelectedFile(event: any) {
        const file = event.target.files[0];
        uploadFile(file);
    }

    const disableDrag = (event: DragEvent) => event.preventDefault();
    function uploadDroppedFile(event: DragEvent) {
        event.preventDefault();
        let files = event.dataTransfer?.files;
        if (files == undefined) return;
        for (let i = 0; i < files.length; i++) {
            let file = files[i];
            let extension = file.name.split(".").at(-1)
            if (extension != "epub") continue;
            uploadFile(file);
        }
    }
</script>


<dialog open>
    <h1> Upload a book </h1>
    <input on:change={uploadSelectedFile} bind:this={fileInput} type="file" style="display: none;">
    <div class="drop" on:click={() => fileInput.click()} on:drop={uploadDroppedFile} on:dragover={disableDrag}>
        <p> Drag your epub files here </p>
        <p> OR </p>
        <p> Click me! </p>
    </div>
</dialog>

<style>
    dialog {
        width: 300px;
        height: 300px;
        text-align: center;
        position: absolute;
    }

    .drop {
        height: 200px;
        background-color: red;
        border: 3px dashed black;
    }
</style>