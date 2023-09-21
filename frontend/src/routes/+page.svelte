<script>
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";

    onMount(() => {
        if (document.cookie == "") {
            goto("/auth");
        }
    });
</script>

<h1> Hello to Page! </h1>