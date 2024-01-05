<template>
    <div id="app_pages">
        <div id="home_page" v-if="page === 'Home Page'"> 
            <div class="centered-div">TIC-TAC-GO</div>            
            <Button @click="startGame" text="Start Game" />
        </div>

        <div v-else-if="page === 'Waiting for Player'">
            <div class="centered-div">Loading...</div>
        </div>

        <div v-else-if="page === 'Start Game'">
            <Board  @game-over="gameOver" :playerNumber="playerNumber"/>
        </div>

        <div v-else-if="page === 'Game Over'">
            <div class="centered-div">{{ playerOutcome }}</div>
            <Button @click="endGame" text="End Game" />
        </div>

        <div v-else>
            <div class="centered-div">404</div>
        </div>
  </div>
</template>

<script>
import Board from './components/Board.vue'
import Button from './components/Button.vue'
import WS from './components/Websocket.vue'

export default {
    // reactive state
    name: 'App',
    components: {
        Board,
        Button,
    },
    data() {
        return{
            page:  "Home Page",
            playerNumber: 0,
            playerOutcome: '',
        }
    },
    mounted() {
        WS.onmessage = (event) => {
            console.log("APP VUE")
            this.page = event.data;
            
            console.log("From app.vue: " + this.page);
            if( this.playerNumber == 0 && this.page == 'Start Game'){
                this.playerNumber = 2
            } else {
                this.playerNumber = 1
            }
        }
    },
    methods: {
        async startGame() {
           WS.send("Join Room")
        },
        gameOver(gameOverMessage) {
            this.page = 'Game Over';
            this.playerOutcome = gameOverMessage;
        },
        async endGame(){
            this.page = 'Home Page';
            WS.send("Terminate Connection")
            WS.onmessage = (event) => {
                console.log("APP VUE")
                this.page = event.data;
                
                console.log("From app.vue: " + this.page);
                if( this.playerNumber == 0 && this.page == 'Start Game'){
                    this.playerNumber = 2
                } else {
                    this.playerNumber = 1
                }
            }
        }
    }
}
</script>

<style scoped>
div #app_pages {    
    height: 100vh;
    width: 100vw;

    position: relative;
    pointer-events: none
}
div #home_page {
    height: 100vh;
    width: 100vw;
}

@font-face {
    font-family: '8bit'; /* Font family name */
    src: url("../public/1up.ttf") /* URL to the font file */
}

.centered-div {
    padding: 10px;
    width: 100%;
    margin: 0 auto 15px; /* 0 top margin, auto left and right margins, 15px bottom margin */

    font-family: '8bit';
    font-size: 150px;
    color: white;
    height: fit-content;
    width: 100%;
    justify-content: center;
    pointer-events: none;
    white-space: nowrap; 
    transform: translate(0%, 50%); 
    display: flex;
    cursor: none;
}
</style>