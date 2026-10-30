# PixSeal Android — GUI nativa (Build106)

App in Go con **Fyne 2.8.2**, per Android **ARM64**. L'elaborazione avviene sul telefono: non serve un PC durante l'uso e il manifesto non richiede accesso a Internet. L'APK è una versione di prova firmata localmente, non una pubblicazione sul Play Store.

## Provala sul telefono

1. Trasferisci `dist/PixSeal-build106-android-arm64.apk` sul telefono e aprilo per installarlo.
2. In **Crea**, scegli una foto JPEG/PNG, inserisci chiave e messaggio, quindi premi **Crea il PNG**. Il messaggio può contenere fino a 64 byte UTF-8, la chiave almeno 8 byte.
3. Salva il PNG. L'app adatta la nuova immagine a un massimo di 1632 pixel sul lato lungo e dimensioni multiple di otto; visualizza e conserva le dimensioni risultanti. La creazione usa forza 48 e profilo automatico, poi verifica il messaggio con una rilettura autenticata.
4. Per una prova digitale, premi **Leggi questo PNG**. Per una prova fisica, stampa il PNG completo con margini bianchi, fotografalo con la fotocamera Android e seleziona la foto nell'app.
5. In **Leggi**, scegli **Foto di una stampa**, usa le dimensioni del PNG **prima della stampa** e la stessa chiave, quindi avvia la lettura. La ricerca può durare diversi minuti: lascia l'app aperta.
6. **Copia messaggio** appare solo dopo autenticazione riuscita. **Salva resoconto** esporta tempi, dimensioni, impronta del file ed esito, senza chiave né testo del messaggio.

La prima versione importa foto già scattate attraverso il selettore file Android; non incorpora una fotocamera. Supporta JPEG e PNG fino a 64 MB / 24 megapixel. La lettura di foto accetta un'immagine canonica fino a 4 megapixel (la preparazione integrata resta entro questo limite). HEIC non è supportato.

## Stato della verifica

- Test del servizio: PNG creato e riletto con autenticazione, chiave errata rifiutata, JPEG decodificato con lo stesso codec della CLI, limiti prima della decodifica, resoconti privi di segreti.
- Test GUI senza display: nessun vecchio messaggio durante una nuova prova, rifiuto di dati non autenticati, blocco di elaborazioni concorrenti, copia del risultato verificato, layout 390×820 e race detector.
- APK: compilazione Android ARM64, integrità ZIP, firma verificata, controllo del manifesto senza permesso Internet. Le licenze dei componenti sono incluse in `assets/NOTICE.txt`.
- **Da verificare su un telefono reale:** avvio, selezione/salvataggio tramite Android, tastiera, comportamento in background, memoria e prove fotografiche. Le anteprime in `dist/gui-preview` sono renderizzazioni dei widget Fyne, non catture da un dispositivo Android.

Il motore di lettura è l'API pubblica Build98 invariata, ma **il target Android non è ancora qualificato fisicamente**. Il caso B/angle rimane irrisolto: questa GUI rende pratiche le prove, non introduce un nuovo algoritmo geometrico.

## Sviluppo e compilazione

Le dipendenze Fyne sono isolate in questo modulo; il modulo Go della CLI non cambia.

```sh
# Dalla radice del repository, con Go 1.26.0:
make gui-test
make gui-run       # anteprima nativa desktop; richiede librerie grafiche
make gui-android   # confeziona e verifica l'APK
```

Strumenti Android: SDK con build-tools 34.0.0, NDK 26.1.10909125 (o percorsi configurati), JDK17, e `go install fyne.io/tools/cmd/fyne@v1.7.4`. Imposta `ANDROID_HOME`, `ANDROID_NDK_HOME` e `JAVA_HOME` se non sono nei percorsi predefiniti. `PIXSEAL_FYNE_TOOL`, `PIXSEAL_APKSIGNER`, `PIXSEAL_ZIPALIGN` e `PIXSEAL_AAPT` consentono percorsi personalizzati.

Build106 allinea `VERSION`, versione interna, `FyneApp.toml` e manifesto Android. Il confezionamento lavora su una copia temporanea, perché Fyne incrementa la configurazione durante la build; controlla anche identità, versione e architettura dell’APK prima di copiarlo in `dist`. `make gui-metadata-check` verifica questi vincoli senza compilare Android.

Il confezionamento usa una **chiave di test locale** in `.tools/gui/android-debug.keystore`, esclusa da Git. Conservala per gli aggiornamenti di prova; non è una chiave di pubblicazione. `PIXSEAL_ANDROID_KEYSTORE` può indicare un'altra chiave di test con alias `pixseal-test`. Il processo allinea l'APK prima della firma e verifica l'artefatto prima di aggiornare il file in `dist`.

Rapporto di verifica: [Build106](../docs/V4_BUILD106_MAINTENANCE.md).

Riferimenti ufficiali: [Fyne mobile packaging](https://docs.fyne.io/started/mobile/), [Fyne main-thread updates](https://docs.fyne.io/started/goroutines/), [Android APK signing](https://developer.android.com/tools/apksigner).
