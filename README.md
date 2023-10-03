# TicTacToe

ecco un esempio di quello che dicevo:
main.go è il file da compilare in WASM e che contiene il codice che deve essere eseguito dal browser
sotto Windows si compila:
```
$env:GOOS="js"; $env:GOARCH="wasm"; go build -o main.wasm .\main.go
```

ed il comando genera  main.wasm che è il file che deve essere caricato dal browser

poi apro un webserver in quella directory e lancio il browser con l'indirizzo del webserver
```
python -m http.server 8000
```

con il browser si va su: 
```
http://127.0.0.1:8000/tictactoe.htm
```

come si vede il file html carica il file main.wasm che è il file che contiene il codice compilato in WASM
ed una volta lanciato esegue:
```
setCode("TESTING");
```

che è la funzione che ho scritto in GO e che viene eseguita dal browser
TESTING sarà un valore ricevuto dal server

quando avviene una vittoria viene concatenato il codice (quello di setCode) con "EOFuserArenaStateread" (una stringa che non da nell'occhio se si decompila il codice) e viene fatto 3 volte lo sha256

alla fine questo dato potrebbe venire inviato al server per verificare che uno abbia effettivamente vinto

In qualsiasi caso un utente potrebbe giocare fintato che non vince e poi inviare il codice al server per verificare che ha vinto... ma questo è un altro discorso