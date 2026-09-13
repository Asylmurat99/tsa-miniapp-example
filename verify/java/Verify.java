import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.MessageDigest;
import java.util.*;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Collectors;

import static java.nio.charset.StandardCharsets.UTF_8;

public class Verify {

    // #region check
    // fields: decoded values; signField: "hash" for the launch context,
    // "sign" for the getPhone envelope.
    static boolean checkSignature(Map<String, String> fields, String secret, String signField) throws Exception {
        var sorted = new TreeMap<>(fields);
        var received = sorted.remove(signField);
        if (received == null) return false;
        var canon = sorted.entrySet().stream()
                .map(e -> e.getKey() + "=" + e.getValue())
                .collect(Collectors.joining("\n"));
        var mac = Mac.getInstance("HmacSHA256");
        mac.init(new SecretKeySpec(secret.getBytes(UTF_8), "HmacSHA256"));
        var calc = HexFormat.of().formatHex(mac.doFinal(canon.getBytes(UTF_8)));
        return MessageDigest.isEqual(calc.getBytes(UTF_8), received.getBytes(UTF_8));
    }
    // #endregion check

    // Minimal reader for vector.json: top-level "secret" and each case's
    // "name", "sign_field" and flat "fields" object. Not a general JSON parser.
    static final Pattern STRING_PAIR = Pattern.compile("\"((?:[^\"\\\\]|\\\\.)*)\"\\s*:\\s*\"((?:[^\"\\\\]|\\\\.)*)\"");

    static String unescape(String s) {
        return s.replace("\\\"", "\"").replace("\\\\", "\\");
    }

    public static void main(String[] args) throws Exception {
        var text = Files.readString(Path.of("../vector.json"));
        var secret = firstValue(text, "secret");
        boolean failed = false;
        for (var block : text.split("\\{\\s*\"name\"")) {
            if (!block.contains("\"fields\"")) continue;
            block = "{\"name\"" + block;
            var name = firstValue(block, "name");
            var signField = firstValue(block, "sign_field");
            var fieldsText = block.substring(block.indexOf("\"fields\"") + 8);
            fieldsText = fieldsText.substring(fieldsText.indexOf('{') + 1);
            // user holds "{...}" inside a string, so the closing brace is
            // found with string state, not with indexOf.
            fieldsText = fieldsText.substring(0, closingBrace(fieldsText));
            var fields = new HashMap<String, String>();
            Matcher m = STRING_PAIR.matcher(fieldsText);
            while (m.find()) fields.put(unescape(m.group(1)), unescape(m.group(2)));
            boolean ok = checkSignature(fields, secret, signField);
            System.out.println(name + ": " + (ok ? "ok" : "FAIL"));
            if (!ok) failed = true;
        }
        System.exit(failed ? 1 : 0);
    }

    static int closingBrace(String s) {
        boolean inString = false;
        for (int i = 0; i < s.length(); i++) {
            char c = s.charAt(i);
            if (inString) {
                if (c == '\\') i++;
                else if (c == '"') inString = false;
            } else if (c == '"') inString = true;
            else if (c == '}') return i;
        }
        throw new IllegalStateException("unterminated fields");
    }

    static String firstValue(String text, String key) {
        Matcher m = Pattern.compile("\"" + key + "\"\\s*:\\s*\"((?:[^\"\\\\]|\\\\.)*)\"").matcher(text);
        if (!m.find()) throw new IllegalStateException(key + " not found");
        return unescape(m.group(1));
    }
}
