export default function initCodes() {
    const form = document.getElementById("remove_code_form");

    form?.addEventListener("change", (e) => {
        const target = e.target as HTMLInputElement;

        if (target.name === "volunteer_code_ids") {
            const codeCheckbox = target.nextElementSibling as HTMLInputElement;
            codeCheckbox.checked = target.checked;
        }
    });
}
