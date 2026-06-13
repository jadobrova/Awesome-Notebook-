<?php
// Блокнот на PHP (запуск: php -S localhost:8000, затем http://localhost:8000/notebook.php?file=test.txt)
$targetFile = $_GET['file'] ?? ($_POST['file'] ?? 'note.txt');
$message = '';

if ($_SERVER['REQUEST_METHOD'] === 'POST' && isset($_POST['content'])) {
    file_put_contents($targetFile, $_POST['content']);
    $message = '✅ Сохранено в ' . htmlspecialchars($targetFile);
}

$content = file_exists($targetFile) ? file_get_contents($targetFile) : '';
$stats = '';
if ($content !== false) {
    $chars = mb_strlen($content);
    $words = str_word_count($content, 0, 'АаБбВвГгДдЕеЁёЖжЗзИиЙйКкЛлМмНнОоПпРрСсТтУуФфХхЦцЧчШшЩщЪъЫыЬьЭэЮюЯя');
    $stats = "Символов: $chars | Слов: $words";
}
?>
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>PHP Блокнот</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 30px; background: #f0f0f0; }
        .container { max-width: 1000px; margin: auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0,0,0,0.1); }
        textarea { width: 100%; height: 60vh; font-family: monospace; font-size: 14px; padding: 10px; }
        input, button { margin: 5px 0; padding: 8px; }
        .status { margin-top: 10px; color: #2c3e50; font-style: italic; }
        .info { background: #e7f3fe; border-left: 4px solid #2196F3; padding: 5px 10px; margin-bottom: 15px; }
    </style>
    <script>
        let timer;
        function autoSave() {
            let content = document.getElementById('editor').value;
            let file = document.getElementById('filename').value;
            if(!file) return;
            fetch(window.location.href, {
                method: 'POST',
                headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                body: new URLSearchParams({content: content, file: file})
            }).then(() => {
                document.getElementById('automsg').innerText = 'Автосохранено в ' + file;
                setTimeout(() => document.getElementById('automsg').innerText = '', 2000);
            });
        }
        document.addEventListener('DOMContentLoaded', () => {
            const editor = document.getElementById('editor');
            editor.addEventListener('input', () => {
                clearTimeout(timer);
                timer = setTimeout(autoSave, 1500);
                let chars = editor.value.length;
                let words = editor.value.trim().split(/\s+/).length;
                document.getElementById('stats').innerText = `Символов: ${chars} | Слов: ${words}`;
            });
        });
    </script>
</head>
<body>
<div class="container">
    <h2>📝 PHP Блокнот с автосохранением</h2>
    <form method="GET" style="margin-bottom: 10px;">
        <input type="text" name="file" id="filename" value="<?= htmlspecialchars($targetFile) ?>" size="40">
        <button type="submit">📂 Загрузить</button>
    </form>
    <form method="POST">
        <textarea name="content" id="editor"><?= htmlspecialchars($content) ?></textarea>
        <input type="hidden" name="file" value="<?= htmlspecialchars($targetFile) ?>">
        <div>
            <button type="submit">💾 Сохранить сейчас</button>
        </div>
    </form>
    <div class="status">
        <span id="stats"><?= $stats ?></span>
        <span id="automsg" style="margin-left: 20px; color: green;"></span>
    </div>
    <?php if($message): ?>
        <div class="info"><?= $message ?></div>
    <?php endif; ?>
</div>
</body>
</html>
