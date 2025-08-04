# DevTUI Shortcut Keys Interface

## Prompt
en #file:README.md quiero implementar que opcionalmente el handler del tipo,#sym:HandlerEdit pueda recibir su valor a traves de un shortcut, ejemplo: 

tengo un manejador llamado tinyWasm que tiene 3 modos de compilación c (coding), d (debug) p (produccion). actualmente si ingreso el valor c en #sym:Change este compilara según lo requerido y asi según el valor que se requiera.. 

entonces seria muy practico que devtui de manera global cuando registra este tipo de manejadores #sym:(*tabSection).AddEditHandler se revisara si cuentan con la firma ShortCuts() []string asi registrarlos de alguna forma en #sym:DevTUI de forma global en una struct especializada en esta feature, asi si el usuario esta en otro tabSection y presiona un shortcut registrado (#sym:(*DevTUI).Update ) la interfaz se traslade al tabSection del handler de quien corresponde ese shortcut y ingrese ese valor al campo Change.. siempre y cuando se este en modo no edición (handleNormalModeKeyboard #sym:(*DevTUI).HandleKeyboard )

antes de llevar llevar a cabo este plan necesito que revises el actual código y en base a tus observaciones elaboremos un documento llamado docs/SHORTCUT_IMPLEMENTATION.md (actualmente vació), en ingles formato prompt con todo lo que debemos realizar detalladamente que archivos/functions/methods/struct/field etc deben ser rectificados/creados etc. debes realizarme las preguntas/sugerencias//recomendaciones pro contras con alternativas debidamente justificadas para tomar las mejores decisiones antes de comenzar con la implementación  

importante: no puedes comenzar la ejecución si no hasta que apruebe este documento