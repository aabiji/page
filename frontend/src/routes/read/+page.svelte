<script lang="ts">
    import BookViewer from "./bookview.svelte";

    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    onMount(() => {
        if (document.cookie == "") {
            goto("/auth");
        }
    });
</script>

<BookViewer bookName="Dune" />
