<script>
const WS = new WebSocket('ws://localhost:8081/public');

WS.onopen = (event) => {
    console.log('WebSocket connection established.');
};

WS.onerror = (error) => {
    console.error('WebSocket error:', error);
};

WS.onclose = () => {
    console.log('WebSocket connection closed');
};
export default WS;
</script>
