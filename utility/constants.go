package utility

/* URL´s Services University*/

//FacultyService - URL query service the information of faculty of a student
const FacultyService = "http://jbpm.udistritaloas.edu.co:8280/services/bienestarProxy/info_institucional/"

//BasicService - URL query service the information basic of a student
const BasicService = "http://jbpm.udistritaloas.edu.co:8280/services/bienestarProxy/info_basica/"

//EnrollmentService -  URL query service the information of value "matricula" of a student
const EnrollmentService = "http://jbpm.udistritaloas.edu.co:8280/services/bienestarProxy/info_recibo/"

//StateService - query URL service the information basic of a student
const StateService = "http://jbpm.udistritaloas.edu.co:8280/services/bienestarProxy/info_estado/"

//AcademicService - query URL service the information basic of a student
const AcademicService = "http://jbpm.udistritaloas.edu.co:8280/services/bienestarProxy/info_academica/"

/* URL's for Files Get and post*/

//FileSavePath - URL of address where the files are saved
const FileSavePath = "C:\\xampp\\htdocs\\Tempfiles\\"

//ServerPath - URL of server where save files
const ServerPath = "http://localhost:80/Tempfiles/"

//RulerPath - MID-API Path - URL of server
const RulerPath = "http://localhost:8090/v1/resultado/"

/* Collections of Data Base*/

//CollectionHistoricFiles - Name of Collection of DB of History of Files
const CollectionHistoricFiles = "apoyoalimentarioarchivos"

//CollectionAdministrator - Name of Collection of DB of Administrator Configurations
const CollectionAdministrator = "apoyoadministracion"

//CollectionGeneral - Name of Collection of DB of Information general students
const CollectionGeneral = "apoyoalimentariogeneral"

//CollectionEconomic - Name of Collection of DB of Information economic students
const CollectionEconomic = "informacioneconomica"
