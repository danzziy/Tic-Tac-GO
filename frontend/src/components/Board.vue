<template>
    <div class="wrapper">
        <div class="board">
            <button v-for="(cell, index) in cells" :key="index" @click="handleClick(index)" :textContent="cell">
                {{ cell }}
            </button>
        </div>
    </div>
</template>

<script>
import WS from './Websocket.vue'

export default {
    // reactive state
    name: 'Board',
    props: {
        playerNumber: Number,
        socket: null
    },
    data() {
        return {
            cells: Array(9).fill(''), // Represents the cells of the board
            gameState: '000000000',
        };
    },
    mounted() {
        if (WS) {
            WS.onmessage = (event) => {
                console.log("From Board.vue: " + event.data);
                var indexOfColon = event.data.indexOf(':');
                var newGameState = event.data.substring(0, indexOfColon);
                this.gameState = newGameState;
                this.cells.forEach((cell, i) => {
                    var playerMove = "X"
                    if(newGameState[i] == 2) {
                        playerMove = "O"
                    }
                    if(newGameState[i] != 0){
                        this.cells[i]= playerMove;
                        console.log("newGameState[i]: " + newGameState[i]);
                    }
                })
                if(!event.data.includes('Ongoing')) {
                    var gameOverMessage = `You ${event.data.substring(indexOfColon + 1)}`;
                    console.log("GAME OVER: " + gameOverMessage)
                    this.$emit('game-over', gameOverMessage);
                } 
            };
        } else {
            console.error('Socket not available');
        }
    },
    methods: {
        async handleClick(index) {
            if (this.cells[index] === '') {
                const message = this.gameState.substring(0, index) + this.playerNumber + this.gameState.substring(index + 1);
                WS.send(message);
            }
        },
    },
    beforeUnmount() {
        // Set WS.onmessage to null to remove the previously assigned function
        WS.onmessage = null;
    },
}
</script>

<style scope>
.wrapper {
    display: flex;
    width: 100%;
    justify-content: center;

}

.board {
    display: inline-grid;
    grid-template-columns: repeat(3, 100px);
    grid-gap: 0px;
    position: relative;
    pointer-events: all;
    justify-content: center;
    align-items: center;
    transform: translate(0%, 20%); 
    background: rgba(0,0,0,0.5);
    text-align: center;
    margin: auto;

}

.board button {
    width: 100px;
    height: 100px;
    background: transparent;
    border: 1px solid whitesmoke;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 2em;
    cursor: pointer;
    color: white;
}
</style>
