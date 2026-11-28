// Examples page: filter by module (color buttons) and by used function.
(function () {
    const cards = Array.from(document.querySelectorAll('#container .mix'));
    const buttons = Array.from(document.querySelectorAll('.legendButton'));
    const input = document.querySelector('#filter_function input');
    const counter = document.getElementById('matches_counter');
    let module = 'all';

    function apply() {
        const q = input.value.trim().toLowerCase();
        let shown = 0;
        for (const card of cards) {
            const fns = card.dataset.fns.toLowerCase().split(' ');
            const okModule = module === 'all' || card.dataset.mod === module;
            const okFn = !q || fns.some(f => f === q || f.includes(q));
            card.classList.toggle('hidden', !(okModule && okFn));
            if (okModule && okFn) shown++;
        }
        counter.textContent = q || module !== 'all' ? shown + ' example' + (shown === 1 ? '' : 's') + ' match' : '';
    }

    for (const b of buttons) {
        b.addEventListener('click', () => {
            module = (module === b.dataset.filter) ? 'all' : b.dataset.filter;
            buttons.forEach(x => x.classList.toggle('active', x.dataset.filter === module));
            apply();
        });
    }

    let timer;
    input.addEventListener('input', () => { clearTimeout(timer); timer = setTimeout(apply, 150); });

    // examples.html?fn=tea.Batch  (links from the cheatsheet)
    const fn = new URLSearchParams(location.search).get('fn');
    if (fn) { input.value = fn; }
    apply();
})();
