<p align="center"><img width="80" src="https://github.com/yukiarrr/SSBird/raw/master/docs/images/logo.png" alt="SSBird logo"></p>
<h2 align="center">SSBird</h2>
<p align="center"><a href="https://github.com/yukiarrr/SSBird/blob/master/README.md">English</a> / <a href="https://github.com/yukiarrr/SSBird/blob/master/README.ja.md">日本語</a></p>

Chrome拡張とSpreadsheetを用いることで、データの作成から反映までを**Chromeだけで**出来るようにしたマスタデータ管理ツールです。

<p align="center"><img width="700" src="https://github.com/yukiarrr/SSBird/raw/master/docs/images/ssbird.gif" alt="SSBird gif"></p>

## 目次

- [概要](#概要)
- [導入手順](#導入手順)
  - [管理者のみ](#管理者のみ)
  - [全メンバー](#全メンバー)
  - [`extension-config.json`について](#extension-configjsonについて)
- [使い方](#使い方)
  - [ダークモードをサポート](#ダークモードをサポート)
  - [各パラメータの説明](#各パラメータの説明)
  - [シートを直接反映する方法](#シートを直接反映する方法)
  - [シート同士をマージして反映する方法](#シート同士をマージして反映する方法)
  - [複数のSpreadsheetのデータを反映する方法](#複数のSpreadsheetのデータを反映する方法)
  - [Pull Requestを作る方法](#Pull-Requestを作る方法)
- [シートのマージ・csv化のルール](#シートのマージcsv化のルール)
  - [基本のルール](#基本のルール)
  - [マージのルール](#マージのルール)

## 概要

マスタデータの管理方法として多いのは、ExcelやSpreadsheetなどでデータを作成し、それをGitHub、そしてDBなどに反映させるやり方だと思われます。  
そのうちSSBirdでは、DBの反映など各プロジェクトの開発環境に依存してしまう部分を除き、データを作成してGitHubにプッシュするまでの工程を負担することで、マスタデータを管理するためのコストを最小限にします。

<p align="center"><img width="500" src="https://github.com/yukiarrr/SSBird/raw/master/docs/images/ssbird-role.jpg" alt="SSBird role"></p>

また、複数人での並行作業やCIでのデータチェックを考慮して、シート同士のマージ機能やPull Requestを作成する機能も備えています。
なお、運用フローとしては、GitHubにcsvがプッシュされたことをトリガーに、CI/CDで各プロジェクトごとに必要なデータの反映することを想定しています。

## 導入手順

### 管理者のみ

1. [Releases](https://github.com/yukiarrr/SSBird/releases)から、使用したいバージョンの`SSBird-vX.X.X.Server.zip`をダウンロードし、解凍してください
2. スクリプトを保存したいGoogle Drive上の場所で右クリックし、`その他 > Google Apps Script`を選択してください
3. コードエディタが開くので、`SSBird/gas/apply.gs`のコードで上書き後に保存してください
4. （省略可能）`apply.gs`の8行目に、`const password = "ここにランダムなパスワードを入力";`のようにパスワードを入力してください（Chrome拡張導入時にApply Passwordとして必要となります）
5. 上のステータスバーから、`リソース > Google の拡張サービス...`を選択し、「Google Sheets API」をONにしてください
6. 上のステータスバーから、`公開 > ウェブ アプリケーションとして導入...`を選択し、「Who has access to the app」でSSBirdを使用するメンバーがアクセスできるように権限を変更してから「Deploy」を押してください
7. 更新後に承認を求められるので、説明に従って承認してください
8. [こちら](#extension-configjsonについて)を参考に、`SSBird/examples/extension-config.json`を編集します（applyUrlはこの手順の5の完了後に表示されるURLを使用してください）
9. 設定ファイルをチームメンバーがアクセスできるGoogle Drive上の場所で右クリックし、`ファイルをアップロード`で編集した`extension-config.json`をアップロードします
10. アップロードした`extension-config.json`を右クリックし、`共有可能なリンクを取得`を選択し、取得した`https://drive.google.com/file/d/XXXXXXX/view?usp=sharing`のうち、`XXXXXXX`の部分をIDとしてメモしてください（Chrome拡張導入時にConfig File Idとして必要となります）

### 全メンバー

1. [Releases](https://github.com/yukiarrr/SSBird/releases)から、使用したいバージョンの`SSBird-vX.X.X.Client.zip`をダウンロードし、解凍してください
2. Macユーザーなら、`SSBird/extension/installer/mac/install.command`を右クリックし、commandキーを押しながら「開く」を選択するとダイアログが出るので、「開く」を選択してください
3. Windowsユーザーなら、`SSBird/extension/installer/windows/install.bat`をダブルクリックしてください
4. 管理者指定のGitHubアカウントでログインしてください
5. [こちら](https://help.github.com/ja/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line)を参考に、参考元の手順の7で「repo」にチェックを入れてから、GitHubのアクセストークンを発行してください
6. `chrome://extensions`にアクセスし、右上のデベロッパーモードをオンにしてから、「パッケージ化されていない拡張機能を読み込む」で`SSBird/extension/app`を選択してください
7. 管理者指定のGoogleアカウントでログインしてください
8. Spreadsheetの画面や、Google DriveのSpreadsheetが存在するフォルダの画面にいくと、Chromeの右上のSSBirdアイコンが押せるようになるので、それを押してください
9. 情報の入力を求められるので、入力してください（GitHub Access Tokenは、この手順の4で取得したものを使用してください）
10. Config File IdとApply Passwordは、管理者から教えてもらってください（省略した場合、Apply Passwordは不要です）

### `extension-config.json`について

```jsonc
{
  // apply.gsのWebアプリケーション公開後に取得できるURL
  "applyUrl": "https://script.google.com/XXXXXXX/exec",

  // プッシュするGitHubリポジトリのURL
  // アクセストークンを発行するユーザーにWrite以上の権限を与えてください
  "repositoryUrl": "https://github.com/XXXXXXX/XXXXXXX",

  // このIDのフォルダからSpreadsheetまでのパスが、GitHub上でのパスとなる
  // フォルダのIDは、そのフォルダをブラウザで開いた状態のURLで
  // https://drive.google.com/drive/u/0/folders/XXXXXXX
  // XXXXXXXの部分が該当する
  "rootFolderId": "XXXXXXX"
}
```

## 使い方

### ダークモードをサポート

|ライト|ダーク|
|:-:|:-:|
|<p align="center"><img width="250" src="https://github.com/yukiarrr/SSBird/raw/master/docs/images/light.jpg" alt="SSBird light mode"></p>|<p align="center"><img width="250" src="https://github.com/yukiarrr/SSBird/raw/master/docs/images/dark.jpg" alt="SSBird dark mode"></p>|

### 各パラメータの説明

|パラメータ|説明|
|:-:|:-:|
|Apply Spreadsheets|Applyの対象となるスプレットシート<br />スプレットシート上では自動選択、Google Drive上では複数選択できるようになる|
|Target Sheet|プッシュの対象となるシート<br />もし存在しなければ、自動で作成される<br />プッシュされるブランチ名はシート名と同じになる|
|Merge Sheets|Target Sheetにマージするシート<br />複数選択可能で、その場合は入力順にマージされていく|
|Commit Message|コミットメッセージ|
|Parent Branch|Target Sheet名のブランチが存在せず、新しく作成される場合の派生元となるブランチ名|
|Create Pull Request|オンの場合、Target Sheetで指定したシートは更新されず、代わりにPull Requestが作成される|

### シートを直接反映する方法

1. Spreadsheetで反映したいシートとデータを作成してください
2. 右上のSSBirdアイコンを押してください
3. 「Target Sheet」で、反映したいシートを選択してください
4. 「Apply」を押してください
5. 「Success 🎉」が出れば成功です

### シート同士をマージして反映する方法

1. Spreadsheetで、反映したいシートとは別に、データを上書きするためのシートを作成してください（上書き用のシートにも、カラムを記述してください）
2. 上書きしたいシートには、上書きしたいデータのみ記述してください
3. 「Target Sheet」で、反映したいシートを選択してください
4. 「Merge Sheets」で、上書きするデータを記述したシートを選択してください（複数選択可で、選択順に上書きしていきます）
4. 「Apply」を押してください
5. 「Success 🎉」が出れば成功です

### 複数のSpreadsheetのデータを反映する方法

1. 反映したいSpreadsheetの入ったGoogle Driveのフォルダ画面に移動してください
2. 「Apply Spreadsheets」で、反映したいSpreadsheetを選択してください
3. 「Target Sheet」「Merge Sheets」に関しては、「[シート同士をマージして反映する方法](#シート同士をマージして反映する方法)」と同じ仕様で、選択したSpreadsheet全てに適応されます
4. 「Apply」を押してください
5. 「Success 🎉」が出れば成功です

### Pull Requestを作る方法

1. 「Create Pull Request」にチェックを入れてください
2. 「Target Sheet」「Merge Sheets」に関しては、「[シート同士をマージして反映する方法](#シート同士をマージして反映する方法)」と同じ仕様です
3. 「Apply」を押してください
4. 「Success 🎉」が出れば成功です（「Target Sheet」がマージ先のブランチ名、「Merge Sheets」がマージするブランチ名になります）

## シートのマージ・csv化のルール

### 基本のルール

- Spreadsheet名がexampleなら、GitHub上ではexample.csvとなります
- ベースとなるシート名がdevelopなら、csvはdevelopブランチにプッシュされます
- カラムは左・上にスペースを空けず記述してください（A1,B1,C1...）
- A列またはカラムが空白なセルは、マージ・csv化どちらにおいても無視されるので、メモをする際にご活用ください（D2にメモするのであれば、A2かD1を空白にする）

<p align="center"><img width="500" src="https://github.com/yukiarrr/SSBird/raw/master/docs/images/ignore-cell.jpg" alt="SSBird ignore cell"></p>

### マージのルール

- A列のセル（A1,A2,A3...）が同じであれば同じデータとみなし上書き、違えば新しいデータとして最下位に追加されます
- データが上書きされる場合、カラムを基準として上書きするので、A列のカラム以外の順番がベースとなるシートと上書きするデータを記述したシートでバラバラでも問題ありません
- 上書きするデータを記述したシートに、ベースとなるシートにない新しいカラムがあれば、新しくカラムが追加されます
