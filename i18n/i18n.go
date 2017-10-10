package i18n // import "ekyu.moe/soda/i18n"

const (
	EN_US = iota
	JA
	ZH_CN
	ZH_TW
)

type Locale int

var (
	YOUR_PUB,
	INPUT_PUB,
	INPUT_PUB_HELP,

	PROMPT_CMD,
	PROMPT_CMD_HELP,
	PROMPT_CMD_ENC,
	PROMPT_CMD_DEC,
	PROMPT_CMD_RAND,
	PROMPT_CMD_CLS,
	PROMPT_CMD_EXIT,
	PROMPT_PLAIN,
	PROMPT_ENCRYPTED,
	PROMPT_OUTPUT,
	PROMPT_OUTPUT_PLAIN,
	PROMPT_OUTPUT_ENCRYPTED,
	PROMPT_OUTPUT_HELP,
	PROMPT_OUTPUT_TERMINAL,
	PROMPT_OUTPUT_EDITOR,

	SESSION_BEGIN,
	ENCRYPTED_BELOW,
	DECRYPT_FAIL,
	PLAIN_BELOW,

	INVALID_PUB,
	INVALID_PLAIN,
	INVALID_ENCRYPTED,

	EXCEPTION_OCCURRED,
	PRESS_ENTER_TO_EXIT string
)

func init() {
	SetLocale(EN_US)
}

func SetLocale(l Locale) {
	switch l {
	case EN_US:
		YOUR_PUB = "Your public key is:"
		INPUT_PUB = "Please input your partner's public key:"
		INPUT_PUB_HELP = `Please copy your public key and send to your partner, while waiting for him
sending his.
Note:
1. Public keys are both one-time, and can only be used in one session.
2. Encryption is end-to-end, but without PKI or TTP, the parter's identity
cannot be verified, i.e. it cannot defend against man-in-the-middle attack.
Therefore, you must authenticate your parter's identity first before any
proceeding.`
		PROMPT_CMD = "What do you want to do next?"
		PROMPT_CMD_HELP = `Both encryption and decryption are designed for text. If you want to encrypt
files, you'd better use some other tools (such as WinRAR), and encrypt the
key via this program before sending to your partner.`
		PROMPT_CMD_ENC = "Encrypt"
		PROMPT_CMD_DEC = "Decrypt"
		PROMPT_CMD_CLS = "Clear the screen"
		PROMPT_CMD_RAND = "Generate a UUIDv4"
		PROMPT_CMD_EXIT = "Exit"
		PROMPT_PLAIN = "Press Enter to launch editor, input plain text, save and quit"
		PROMPT_ENCRYPTED = "Press Enter to launch editor, input encrypted text, save and quit"
		PROMPT_OUTPUT_PLAIN = "Use editor or terminal to print the plain text? (editor is preferred)"
		PROMPT_OUTPUT_ENCRYPTED = "Use editor or terminal to print the encrypted text? (editor is preferred)"
		PROMPT_OUTPUT_HELP = `If editor is selected, a temp file will be created to store the content, and it
will be destroyed after closing the editor, while printing to terminal is very
likely to retain a history. That's why editor is preferred.`
		PROMPT_OUTPUT_TERMINAL = "Terminal"
		PROMPT_OUTPUT_EDITOR = "Editor"
		SESSION_BEGIN = "Session Begin"
		ENCRYPTED_BELOW = "The encrypted text is as follows. You can copy and send it to your partner."
		DECRYPT_FAIL = "Failed to decrypt. The encrypted text may have been compromised."
		PLAIN_BELOW = "The plain text is as follows."
		INVALID_PUB = "Bad public key."
		INVALID_PLAIN = "Plain text should not be empty."
		INVALID_ENCRYPTED = "Bad encrypted text."
		EXCEPTION_OCCURRED = "An error has occurred."
		PRESS_ENTER_TO_EXIT = "Press Enter to exit..."
	case JA:
		YOUR_PUB = "あなたの公開鍵は："
		INPUT_PUB = "相手の公開鍵を入力してください："
		INPUT_PUB_HELP = `
あなたの公開鍵をコピーして、相手に送信すると、相手の公開鍵を受けるうちにしばらく
お待ちください。
お知らせ：
1. 両方の公開鍵はワンタイムであり、一度しか使用できないようになっています。
2. 暗号化はエンドツーエンドですが、公開鍵基盤はない場合、中間者攻撃は防止できな
い。よって、相手の身分認証は必要です。`
		PROMPT_CMD = "次に何をすべきか？"
		PROMPT_CMD_HELP = "暗号化と復号化は文字列に対します。ファイルを暗号化の場合、他のツール(WinRAR\nとか)を使って、鍵をここに暗号化して相手に送信してください。"
		PROMPT_CMD_ENC = "暗号化"
		PROMPT_CMD_DEC = "復号化"
		PROMPT_CMD_CLS = "ターミナルをクリア"
		PROMPT_CMD_RAND = "UUIDv4 を生成"
		PROMPT_CMD_EXIT = "終了"
		PROMPT_PLAIN = "Enter キーを押してエディタを開いて、プレーンテキストを入力して、セーブして\nエディタを終了してください"
		PROMPT_ENCRYPTED = "Enter キーを押してエディタを開いて、暗号化された テキストを入力して、セーブして\nエディタを終了してください"
		PROMPT_OUTPUT_PLAIN = "出力の方、エディタまたはターミナルを使いますか？（エディタに推薦）"
		PROMPT_OUTPUT_ENCRYPTED = "出力の方、エディタまたはターミナルを使いますか？（エディタに推薦）"
		PROMPT_OUTPUT_HELP = "エディタを選択した場合、一時ファイルが作成されて、エディタが終了する時に削除\nされます。ターミナルを選択した場合、履歴ログの保持の可能性は高いです。"
		PROMPT_OUTPUT_TERMINAL = "ターミナル"
		PROMPT_OUTPUT_EDITOR = "エディタ"
		SESSION_BEGIN = "セッション開始"
		ENCRYPTED_BELOW = "暗号化されたテキストは以下の通りです"
		DECRYPT_FAIL = "復号化失敗しました、テキストは危殆化されたの可能性は高いです。"
		PLAIN_BELOW = "復号化されたテキストは以下の通りです"
		INVALID_PUB = "無効の公開鍵。"
		INVALID_PLAIN = "空のテキストが無効です。"
		INVALID_ENCRYPTED = "無効のテキスト。"
		EXCEPTION_OCCURRED = "例外が発生しました"
		PRESS_ENTER_TO_EXIT = "Enter を押して終了します..."
	case ZH_TW:
		YOUR_PUB = "你的公鑰為："
		INPUT_PUB = "請輸入對方的公鑰："
		INPUT_PUB_HELP = `請將你的公鑰複製下來，發送給對方，同時等待對方發來他的公開金鑰。
注意：
1.雙方的公鑰都是一次性的，只能用於當次對談。
2.加密是端到端的，但在沒有公鑰基礎設施或可信任的第三方輔助的情况下，對方
的身份是無法保證的，即在中間人攻擊面前是極其脆弱的。所以在操作前請務必
先確認對方的身份。`
		PROMPT_CMD = "請問接下來要做什麼？"
		PROMPT_CMD_HELP = "加密與解密針對的都是文字，如果要加密檔案，可以採取別的對稱加密手段\n（如使用 WinRAR 加密），然後再將金鑰通過這裡加密後發給對方。"
		PROMPT_CMD_ENC = "加密"
		PROMPT_CMD_DEC = "解密"
		PROMPT_CMD_CLS = "清屏"
		PROMPT_CMD_RAND = "生成一段 UUIDv4"
		PROMPT_CMD_EXIT = "退出"
		PROMPT_PLAIN = "按回車鍵打開編輯器，輸入明文，然後保存並關閉編輯器"
		PROMPT_ENCRYPTED = "按回車鍵打開編輯器，輸入密文，然後保存並關閉編輯器"
		PROMPT_OUTPUT_PLAIN = "使用編輯器顯示明文還是直接輸出明文到終端？（推薦選擇編輯器）"
		PROMPT_OUTPUT_ENCRYPTED = "使用編輯器顯示密文還是直接輸出密文到終端？（推薦選擇編輯器）"
		PROMPT_OUTPUT_HELP = "如選擇編輯器，則會創建一個暫存檔案來存放內容，在關閉編輯器後會被自動删除，\n而输出到終端的話很可能會留下歷史。所以推薦以編輯器的方式查看。"
		PROMPT_OUTPUT_TERMINAL = "終端"
		PROMPT_OUTPUT_EDITOR = "編輯器"
		SESSION_BEGIN = "對談開始"
		ENCRYPTED_BELOW = "以下為剛剛加密的密文，可以複製下來發送給對方了。"
		DECRYPT_FAIL = "解密失敗。密文很可能已被污染。"
		PLAIN_BELOW = "以下為剛剛解密的明文"
		INVALID_PUB = "公開金鑰格式不正確。"
		INVALID_PLAIN = "明文不能為空。"
		INVALID_ENCRYPTED = "密文格式不正確。"
		EXCEPTION_OCCURRED = "發生异常"
		PRESS_ENTER_TO_EXIT = "按回車鍵退出…"
	case ZH_CN:
		YOUR_PUB = "你的公钥为："
		INPUT_PUB = "请输入对方的公钥："
		INPUT_PUB_HELP = `请将你的公钥复制下来，发送给对方，同时等待对方发来他的公钥。
注意：
1. 双方的公钥都是一次性的，只能用于当次会话。
2. 加密是端到端的，但在没有公钥基础设施或可信任的第三方辅助的情况下，对方
的身份是无法保证的，即在中间人攻击面前是极其脆弱的。所以在操作前请务必先
确认对方的身份。`
		PROMPT_CMD = "请问接下来要做什么？"
		PROMPT_CMD_HELP = "加密与解密针对的都是文本，如果要加密文件，可以采取别的对称加密手段（如使用\nWinRAR 加密），然后再将密钥通过这里加密后发给对方。"
		PROMPT_CMD_ENC = "加密"
		PROMPT_CMD_DEC = "解密"
		PROMPT_CMD_CLS = "清屏"
		PROMPT_CMD_RAND = "生成一段 UUIDv4"
		PROMPT_CMD_EXIT = "退出"
		PROMPT_PLAIN = "按回车键打开编辑器，输入明文，然后保存并关闭编辑器"
		PROMPT_ENCRYPTED = "按回车键打开编辑器，输入密文，然后保存并关闭编辑器"
		PROMPT_OUTPUT_PLAIN = "使用编辑器显示明文还是直接输出明文到控制台？（推荐选择编辑器）"
		PROMPT_OUTPUT_ENCRYPTED = "使用编辑器显示密文还是直接输出密文到控制台？（推荐选择编辑器）"
		PROMPT_OUTPUT_HELP = "如选择编辑器，则会创建一个临时文件来存放内容，在关闭编辑器后会被自动删除，\n而打印到终端的话很可能会留下历史。所以推荐以编辑器的方式查看。"
		PROMPT_OUTPUT_TERMINAL = "控制台"
		PROMPT_OUTPUT_EDITOR = "编辑器"
		SESSION_BEGIN = "会话开始"
		ENCRYPTED_BELOW = "以下为刚刚加密的密文，可以复制下来发送给对方了"
		DECRYPT_FAIL = "解密失败。密文很可能已被污染。"
		PLAIN_BELOW = "以下为刚刚解密的明文"
		INVALID_PUB = "公钥格式不正确。"
		INVALID_PLAIN = "明文不能为空。"
		INVALID_ENCRYPTED = "密文格式不正确。"
		EXCEPTION_OCCURRED = "发生异常"
		PRESS_ENTER_TO_EXIT = "按 Enter 键退出..."
	}
}
