# apoyoalimentario_CRUD_API

Este proyecto fue generado con Framework BeeGo version 1.9.1., base de datos MongoBD version > 3.6 y Servidor Apache

Esta aplicación sirve para realizar las peticiones entre cliente y Base de Datos de la inscripción al Apoyo Alimentario de la Universidad Distrital Francisco José de Caldas

# IMPORTANTE
Es necesario que se tenga el CLIENTE y MID_API

Configuración de puerto<br>
Modifique la propiedad llamada httpport del archivo app.conf que está ubicado en conf/httpport<br>

Configuración de Base de datos MongoDB<br><br>
1 Modifique las propiedades llamadas <br>
    - mongo_host : con el host<br>
    - mongo_db   : con el nombre de la base de datos<br>
    - mongo_user : usuario<br>
    - mongo_pass : contraseña<br>
del archivo app.conf que está ubicado en conf/httpport

2 en el archivo utility/contants.go<br>
Se encuentran las colecciones que debe tener la base de datos<br>
-apoyoalimentarioarchivos<br>
-apoyoadministracion<br>
-apoyoalimentariogeneral<br>
-informacioneconomica<br>
<br>
<br>
#Configurar servidor apache
en el archivo utility/contants.go<br>
 constante llamada ServerPath por defecto http://localhost:80/Tempfiles/

# Para desplegar
Debe ejecutar una vez instalado Go y Beego, el comando bee run para ejecutarlo en un servidor local, la aplicación corre por defecto en http://localhost:8086/.

Si desea cambiarla, modifique la propiedad httpport del archivo app.conf esta ubicado en la carpeta conf del proyecto.
