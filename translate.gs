function doGet(e) {
    var p = e.parameter;
    var translatedText = LanguageApp.translate(p.text, p.source, p.target);
    var body;
    if (translatedText) {
        body = {
          code: 200,
          text: translatedText
        };
    } else {
        body = {
          code: 400,
          text: "Bad Request"
        };
    }

    var response = ContentService.createTextOutput();
    response.setMimeType(ContentService.MimeType.JSON);
    response.setContent(JSON.stringify(body));

    return response;
}
