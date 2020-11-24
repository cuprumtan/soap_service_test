from googletrans import Translator
from spyne import Application, rpc, ServiceBase, Unicode
from lxml import etree
from spyne.protocol.soap import Soap11
from spyne.server.wsgi import WsgiApplication
import yaml


def load_config():
    with open('config.yml') as file:
        config = yaml.load(file, Loader=yaml.FullLoader)
    return config


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
    config = load_config()
    from wsgiref.simple_server import make_server
    server = make_server(config['soap_host'], int(config['soap_port']), application)
    server.serve_forever()
