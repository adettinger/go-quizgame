import { WebSocketControl } from "./WebSocketControl";


export function GameHost () {
  return (
    <WebSocketControl />
  );
};


// export function GameHost() {

    
//     const handleStartGame = () => {
//         const socket = new WebSocket('ws://localhost:8080/ws');
//         socket.addEventListener('open', (event) => {
//             console.log('Connected to WebSocket server');
//             socket.send('Hello Server!');
//         });
//     };

//     const handleEndGame = () => {
//         socket.close();
//         socket.addEventListener('close', (event) => {
//             console.log('Socket closed');
//         });
//     };


//     return (
//         <Flex align="center" justify="center">
//             <Button
//                 onClick={handleStartGame}
//             >
//                 Start Game
//             </Button>

//             <Button
//                 onClick={handleEndGame}
//             >
//                 End Game
//             </Button>
//         </Flex>
//     )
// }