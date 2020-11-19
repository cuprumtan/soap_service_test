from googletrans import Translator
from spyne import Application, rpc, ServiceBase, Unicode
from lxml import etree
from spyne.protocol.soap import Soap11
from spyne.protocol.json import JsonDocument
from spyne.server.wsgi import WsgiApplication


class SOAP(ServiceBase):
    @rpc(Unicode, Unicode, Unicode, _returns=Unicode)
    def GoogleTranslate(ctx, source_text, source_language, destination_language):
        translator = Translator()
        print(etree.tostring(ctx.in_document))
        service_answer = ""
        # костыль для баги в googletrans
        while service_answer == "":
            try:
                result = translator.translate(str(source_text), src=source_language, dest=destination_language)
                service_answer = result.text
            except Exception:
                translator = Translator()

        return service_answer


app = Application([SOAP], tns='Translator',
                  in_protocol=Soap11(validator='lxml'),
                  out_protocol=Soap11())
application = WsgiApplication(app)
if __name__ == '__main__':
    from wsgiref.simple_server import make_server
    server = make_server('0.0.0.0', 8000, application)
    server.serve_forever()
