package utility

/* URL´s de Servicios de la Universidad*/

//FacultyService - URL retorna servicio con la informacion institucional de un estudiante
const FacultyService = "http://jbpm.udistritaloas.edu.co:8280/services/bienestarProxy/info_institucional/"

//BasicService - URL retorna servicio con la informacion basica de un estudiante
const BasicService = "http://jbpm.udistritaloas.edu.co:8280/services/bienestarProxy/info_basica/"

//EnrollmentService -  URL retorna servicio con la informacion de la matricula sin sitematizacion de un estudiante
const EnrollmentService = "http://jbpm.udistritaloas.edu.co:8280/services/bienestarProxy/info_recibo/"

//StateService - URL retorna servicio con la informacion del estado ACTIVO o INCTIVO de un estudiante
const StateService = "http://jbpm.udistritaloas.edu.co:8280/services/bienestarProxy/info_estado/"

//AcademicService - URL retorna servicio con la informacion academica de un estudiante
const AcademicService = "http://jbpm.udistritaloas.edu.co:8280/services/bienestarProxy/info_academica/"

/* URL's*/

//FileSavePath - URL de direccion Local donde los archivos se van a guardar
const FileSavePath = "C:\\xampp\\htdocs\\Tempfiles\\"

//ServerPath - URL del servidor donde los archivos son llamados APACHE
const ServerPath = "http://localhost:80/Tempfiles/"

/*Collecciones de Base de Datos MONGODB*/

//CollectionHistoricFiles - Nombre de la coleccion de BD del Historico de archivos
const CollectionHistoricFiles = "apoyoalimentarioarchivos"

//CollectionAdministrator - Nombre de la coleccion de BD de la configuracion del Administrador
const CollectionAdministrator = "apoyoadministracion"

//CollectionGeneral - Nombre de la coleccion de BD de la Informacion general de los estudiantes
const CollectionGeneral = "apoyoalimentariogeneral"

//CollectionEconomic - Nombre de la coleccion de BD de la Informacion economica de los estudiantes
const CollectionEconomic = "informacioneconomica"
