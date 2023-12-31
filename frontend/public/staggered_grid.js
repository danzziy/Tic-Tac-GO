
let columns = Math.floor(document.body.clientWidth / 50),
    rows = Math.floor(document.body.clientHeight / 50)

const wrapper = document.getElementById("tiles");


var count = -1
const handleOnClick = index => {
    count++
   
    anime({
        targets: ".tile",
        delay: anime.stagger(50, {
            grid: [columns, rows],
            from: index
        }),
        scale: [{ value: [1, 0.9], duration: 500}, { value: 1, duration: 1000, easing: 'easeInOutElastic(1, .9)'}],
    })
}

const createTile = index => {
    const tile = document.createElement('div');

    tile.classList.add("tile");

    tile.onclick = e => handleOnClick(index);
    
    return tile
}

const createTiles = quantity => {
    Array.from(Array(quantity)).map((tile, index) => {
        wrapper.appendChild(createTile(index));
    })
}

createTiles(columns * rows);

const createGrid = () => {
    wrapper.innerHTML = "";
    
    columns = Math.floor(document.body.clientWidth / 50);
    rows = Math.floor(document.body.clientHeight / 50);

    wrapper.style.setProperty("--columns", columns);
    wrapper.style.setProperty("--rows", rows);

    createTiles(columns * rows);
}

createGrid();

window.onresize = () => createGrid();

const body = document.body;

onmousemove = (e) => {
    let {x, y} = e;
    body.style.setProperty("--x", `${x}px`);
    body.style.setProperty("--y", `${y}px`);
}
