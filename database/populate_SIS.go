package database

import (
	"log"
)

func SeedUserTypes() error {
	if db == nil {
		log.Fatal("Database not initialized. Call SetDB first.")
	}

	records := []struct {
		id      int
		descrAr string
		descrEn string
	}{
		{2, "طالب", "Student"},
		{8, "مسؤل المتقدمين", "Admission Admin"},
		{10, "مدير السرية والامان", "Security Admin"},
		{12, "كفيل", "Sponsor"},
		{13, "ادارة الأمن", "Security Department"},
		{14, "عميد شوؤن الطلاب", "Dean of Student Affairs"},
		{15, "رئيس قسم", "Head of Department"},
		{16, "عضو هيئة التدريس", "College Staff"},
		{17, "نائب رئيس الجامعة للشئون الاكاديمية", "Vice president for Academic Affairs"},
		{18, "نائب رئيس الجامعة للشؤون المالية والاداريه", "Vice president for Financial Administrative Affairs"},
		{19, "نائب رئيس الجامعة لخدمة المجتمع", "Vice president for Community Service"},
		{20, "الخريجين", "Alumni"},
		{21, "ولي الامر", "Parents"},
		{22, "رئيس الجامعة", "University President"},
		{24, "عميد الكلية", "College Dean"},
		{25, "مسئول المصروفات الدراسية", "Fee Affairs Admin"},
		{26, "مدير  النظام", "System Admin"},
		{27, "مدير  التسجيل", "Registration Admin"},
		{28, "موظف", "Employee"},
		{29, "أدارة المبانى والصيانة", "Buildings and Maintenance Administration"},
		{30, "أدارة بيع الكتب الدراسية", "Bookshop Admin"},
		{31, "Internal Auditor Module", "Internal Auditor Module"},
		{32, "جهات خارجية", "External Parties"},
		{33, "Technical Support", "Technical Support"},
		{34, "Library Information Services", "Library Information Services"},
		{35, "DeanShip of PG studies ", "DeanShip of PG studies"},
		{36, "نائب رئيس الجامعة لشؤون الشراكات والتنمية", "Vice president for Partnerships and Development"},
	}

	stmt, err := db.Prepare(`
		INSERT OR IGNORE INTO SE_CODE_USER_TYPE (SE_CODE_USER_TYPE_ID, DESCR_AR, DESCR_EN)
		VALUES (?, ?, ?);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, r := range records {
		if _, err := stmt.Exec(r.id, r.descrAr, r.descrEn); err != nil {
			return err
		}
	}

	log.Println("✅ SE_CODE_USER_TYPE table populated successfully")
	return nil
}

func SeedAccounts() error {
	if db == nil {
		log.Fatal("Database not initialized. Call SetDB first.")
	}

	type Account struct {
		ID       int
		DescAr   string
		DescEn   string
		UserType int
	}

	records := []Account{
		{2, "طالب", "Student", 2},
		{3, "حساب الطالب", "Student Module", 2},
		{4, "حساب عميد شئون الطلاب", "Dean of Student Affairs Module", 14},
		{5, "حساب رئيس الجامعة", "University President Module", 22},
		{6, "حساب نائب رئيس الجامعة للشئون الاكاديمية", "Vice president for Academic Affairs Module", 17},
		{7, "حساب مسئول المتقدمين", "Admission Admin Module", 8},
		{8, "مسئول المتقدمين", "Admission Admin", 8},
		{9, "حساب مسئول المصروفات الدراسية", "Fees Admin Module", 25},
		{10, "قسم السرية والامان", "Security Admin Module", 10},
		{11, "مدير السرية والأمان", "Security Admin", 10},
		{12, "حساب مشرف نظام", "System Admin Module", 26},
		{13, "حساب كفيل", "Sponsor Module", 12},
		{14, "كفيل", "Sponsor", 12},
		{15, "قسم ادارة الامن بالجامعة", "Security Department Module", 13},
		{16, "عميد شئون الطلاب", "Dean of Student Affairs", 14},
		{17, "رئيس قسم", "Head of Department", 15},
		{18, "عضو هيئة التدريس", "College Staff Member", 16},
		{19, "حساب الخريجين", "Alumni Module", 20},
		{20, "حساب رئيس قسم", "Head of Department Module", 15},
		{21, "حساب عضو هيئة التدريس", "College Staff Member Module", 16},
		{22, "حساب ولي الامر", "Parents Module", 21},
		{24, "حساب عميد الكلية", "College Dean Module", 24},
		{25, "نائب رئيس الجامعة للشئون الاكاديمية", "Vice president for Academic Affairs", 17},
		{26, "عميد الكلية", "College Dean", 24},
		{27, "مسئول المصروفات الدراسية", "Fees Admin", 25},
		{28, "مشرف نظام", "System Admin", 26},
		{29, "مدير التسجيل", "Registration Admin", 27},
		{30, "حساب مدير التسجيل", "Registration Admin Module", 27},
		{31, "مسئول الجداول الدراسية", "Schedule Affairs", 15},
		{32, "منسق الدراسات العليا", "PG Co-ordinator", 15},
		{33, "مسئول التدريب", "Training Affairs", 15},
		{73, "مجموعة صيانة بيانات اللائحة", "Bylaw Maintenance Group", 26},
		{74, "سكرتيرة القبول وشئون الخريجين", "Admission Secretary", 8},
		{75, "رئيس القبول", "Admission Chief", 8},
		{76, "مشرف القبول وشئون الخريجين", "Admission Supervisor", 8},
		{77, "كاتب أول القبول وشئون الخريجين", "Admission First Writer", 8},
		{78, "فني سجلات القبول وشئون الخريجين", "Admission Technician Records", 8},
		{79, "كاتب القبول وشئون الخريجين", "Admission Writer", 8},
		{80, "رئيس شعبة التسجيل والإرشاد الأكاديمي", "Registration Chief Division", 27},
		{81, "أخصائي التسجيل والإرشاد الأكاديمي", "Registration Specialist", 27},
		{82, "موظفي التسجيل", "Registration Staff", 27},
		{83, "مشرف الجداول الدراسية للتسجيل والإرشاد الأكاديمي", "Registration Course Schedule Supervisor", 27},
		{84, "فني التسجيل والإرشاد الأكاديمي", "Registration Technician", 27},
		{85, "كاتب أول التسجيل والإرشاد الأكاديمي", "Registration First Writer", 27},
		{86, "كاتب التسجيل والإرشاد الأكاديمي", "Registration Writer", 27},
		{87, "الخدمات الطلابية", "Students Service", 14},
		{88, "رئيس الاعفاءات", "Exemptions Head", 14},
		{89, "دائرة التوجية والارشاد", "Advising and  Guidance Department", 14},
		{90, "موظف أمن", "Security Employee", 13},
		{91, "شعبة التوجية الاجتماعى", "Social Guidance section", 14},
		{92, "شعبة الارشاد النفسى", "Psychological Counseling", 14},
		{94, "عميد القبول والتسجيل - التسجيل", "Dean of Admission & Registration - Registeration", 27},
		{95, "دائرة نظم المعلومات وتحليل البيانات والإحصاءات - التسجيل", "Department of Information Data Analysis and Statistics - Registration", 27},
		{96, "عميد القبول والتسجيل -  القبول", "Dean of Admission & Registration - Admission", 8},
		{97, "دائرة نظم المعلومات وتحليل البيانات والإحصاءات  - القبول", "Department of Information Data Analysis and Statistics  - Admission", 8},
		{98, "RegisterAddDrop", "RegisterAddDrop", 27},
		{99, "سكرتير العميد", "Dean Secretary", 27},
		{100, "فنى جداول", "Schedule Technician", 27},
		{101, "مكتب الرئيس", "President Office", 27},
		{102, "حساب نائب رئيس الجامعة لشئون الطلاب ", "vice president for academic affairs  account", 17},
		{103, " رئيس الجامعة", "President ", 22},
		{104, "مشرف جداول كلية البحرين للمعلمين", "BTC Scheduling Supervisor", 27},
		{105, "مدير الخدمات الطلابية", "Students Service Admin", 14},
		{106, "Assistant", "Assistant", 24},
		{107, "Student Reg Complaints", "Students Reg Complaints", 27},
		{108, "مدير ضمان الجودة", "QAAC Director", 17},
		{109, "موظف ضمان الجودة", "QAAC Employee", 17},
		{110, "Schedule Reports", "Schedule Reports", 27},
		{111, "Public Relations - Reg", "Public Relations - Reg", 27},
		{112, "مساعد رئيس القسم", "Head of Department Assistant", 15},
		{113, "Quality Office Director", "Quality Office Director", 17},
		{114, "مدقق داخلي - Not used", "Internal Auditor - Not used", 15},
		{115, "Internal Auditor", "Internal Auditor", 27},
		{117, "عميد الدراسات العليا والبحث العلمي", "Dean of PG studies ", 35},
		{118, "مدير مكتب الرئيس", "President Office Admin", 22},
		{119, "نظام لانتخابات للطلاب بجامعة البحرين", "University of Bahrain Voting System", 13},
		{120, "موظفين الاعفاءات", "Exemptions Employee", 14},
		{121, "ادارة الامن بجامعة البحرين", "UOB Security Department", 13},
		{122, "موظف", "Employee", 28},
		{123, "كاشير", "Cheque Cashier", 25},
		{124, "أمين صندوق", "Cashier", 25},
		{125, "الشيكات", "Cheques", 25},
		{126, "محاسب", "Accountant ", 25},
		{128, "Applicant Interview", "Applicant Interview", 8},
		{129, "أدارة المبانى والصيانة", "Buildings and Maintenance Administration", 29},
		{130, "Gulf CX", "Gulf CX", 8},
		{131, "الخريجين", "Alumni", 20},
		{132, " ولي الامر", "Parents", 21},
		{133, "Pressident Advisor", "Pressident Advisor", 22},
		{134, "Reset Student Password Account", "Reset Student Password Account", 27},
		{135, "مدير بيع الكتب الدراسية", "Bookshop Admin", 30},
		{136, "BookShop Employee", "BookShop Employee", 30},
		{137, "BookShop Cashier", "BookShop Cashier", 30},
		{138, "BookShop Seniors", "BookShop Seniors", 30},
		{139, "مدقق داخلي - مالي", "Financial Internal Auditor", 25},
		{140, "BookShop Module", "BookShop Module", 30},
		{141, "Records & Graduates Affairs Chief", "Records & Graduates Affairs Chief", 8},
		{142, "Transfer Supervisior", "Transfer Supervisior", 8},
		{143, "كاتب اول", "First Clerk", 2},
		{144, "نائب الرئيس لخدمات تقنية المعلومات والشؤون الإدارية والمالية", "Vice President for IT , Admin. and Finance", 18},
		{145, "نائب الرئيس لتقنية المعلومات والشؤون الإدارية والمالية", "Vice President for IT , Admin and Finance", 18},
		{146, "منسق ضمان الجودة", "QAAC Coordinator", 17},
		{147, "مدير مركز اللغة الانجليزية", "Director of English Language Center", 27},
		{148, "استقبال مقابلات الدراسات العليا", "PG Interviews Reception", 8},
		{149, "رئيس شعبة القبول", "Head of Admission Section", 8},
		{150, "رئيس شعبة السجلات", "Head of Records Section", 8},
		{151, "مدير دائرة القبول", "Director of Admission Dept.", 8},
		{152, "مشرف وحدة القبول", "Supervisor of Admission Unit", 8},
		{153, "مشرف وحدة التحويل", "Supervisor of Transfer Unit", 8},
		{154, "مشرف وحدة الشهادات ومتابعة الخريجين ", "Supervisor of Certificates Unit & Graduates", 8},
		{155, "Demonstrator", "Demonstrator", 16},
		{156, "Internal Auditor Module", "Internal Auditor Module", 31},
		{157, "حساب المدقق الداخلى", "internal Auditor Account", 31},
		{158, "مدير دائرة التوجيه والإرشاد بعمادة شؤون الطلبة", "Director of the Guidance and Counseling Department of the Deanship of Student Affairs", 14},
		{159, "حساب الجهات الخارجية", "External Parties Account", 32},
		{160, "موظفى جسر الملك فهد", "King Fahad Causeway Employees", 32},
		{161, "Black Board Team", "Black Board Team", 27},
		{162, "Teams and Blackboard Admin", "Teams and Blackboard Admin", 33},
		{163, "مدير الموارد البشرية", "HR Manager", 18},
		{164, "موظف موارد بشرية", "HR Employee", 18},
		{166, "Library and information Templete", "Library and information Templete", 34},
		{167, "Library and information ", "Library and information ", 34},
		{168, "مستخدمي الإعفاءات", "Exemptions users", 14},
		{169, "BTC Scheduling Assistant", "BTC Scheduling Assistant", 27},
		{170, "BTC Scheduling Assistant1", "BTC Scheduling Assistant1", 27},
		{171, "متابعة طلبة الدراسات العليا وتسجيل مقرر الأطروحة", "Follow-up of PG Students & Registering Thesis Course", 16},
		{172, "PG Studies", "PG Studies", 35},
		{173, "QAAC Supervisor", "QAAC Supervisor", 17},
		{174, "لجنة الحالة الطلابية", "Case Study Committee Head", 17},
		{175, "View Student Data", "View Student Data", 27},
		{176, "خدمات الطالب", "Student Service", 32},
		{177, "Alumni Club Manager", "Alumni Club Manager", 14},
		{178, "New Accountant", "New Accountant", 25},
		{179, "مساعد عميد الدراسات العليا", "PG Dean - Assistant", 8},
		{180, "دراسات عليا - مساعد", "PG Studies - Assistant", 35},
		{181, "منسق الدراسات العليا - مساعد", "PG Co-ordinator - Assistant", 15},
		{182, "حساب إضافة فصول قديمة", "Add Old Semester Account", 27},
		{183, "معيد", "TEACHING ASSISTANT", 16},
		{184, "مشرف اللوائح الجامعية و إعدادات النظام", "Bylaw and setup Supervisor ", 26},
		{185, "حساب عميد الدراسات العليا والبحث العلمي", "Dean of PG Studies Module", 35},
		{186, "متقدمي كلية البحرين للمعلمين", "BTC Admission", 8},
		{187, "راعي شؤون الطلاب", "Student Affairs Sponsor", 14},
		{188, "إدارة اللوائح الأكاديمية", "Bylaw Management", 26},
		{189, "تقرير العبء الأكاديمي لأعضاء هيئة التدريس", "Faculty Academic Load Report Viewer", 28},
		{190, "Academic Faculty Load - deleted", "Academic Faculty Load -deleted", 15},
		{191, "Academic Faculty Load", "Academic Faculty Load", 28},
		{192, "VPA Coordinator", "VPA Coordinator", 17},
		{193, "QAAC Coordinator Other", "QAAC Coordinator Other", 17},
		{194, "Student Case Committee View", "Student Case Committee View", 17},
		{195, "لجنة الحالة الطلابية عضو", "Case Study Committee Member", 17},
		{196, "الحضور الموارد البشرية", "Attendance HR", 27},
		{197, "مدير الحضور و الانصراف", "Attendance Admin", 27},
		{198, "Technician", "Technician", 28},
		{199, "لجنة الدراسات العليا", "PG Committee", 35},
		{200, "تسليم الكتب ", "Book Delivery", 30},
		{201, "مساعد مدير مركز اللغة الانجليزية", "Assistant Director of English Language Center", 27},
		{202, "إدارة كشوف الدرجات", "Manage Trascript", 27},
		{203, "Refund Account", "Refund Account", 25},
		{204, "مدير الخريجين", "Alumni Manager", 36},
		{205, "نائب الرئيس للشراكات والتنمية", "VP for Partnerships and Development", 36},
		{206, "Training Officer", "Training Officer", 36},
	}

	stmt, err := db.Prepare(`
		INSERT OR IGNORE INTO SE_ACCNT (SE_ACCNT_ID, DESCR_AR, DESCR_EN, SE_CODE_USER_TYPE_ID)
		VALUES (?, ?, ?, ?);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, r := range records {
		if _, err := stmt.Exec(r.ID, r.DescAr, r.DescEn, r.UserType); err != nil {
			return err
		}
	}

	log.Println("✅ SE_ACCNT table populated successfully")
	return nil
}

func SeedEdCodeStatusCat() error {
	if db == nil {
		log.Fatal("Database not initialized. Call SetDB first.")
	}

	type Cat struct {
		ID      int
		DescrEn string
		DescrAr string
	}

	records := []Cat{
		{1, "Cheque STATUS", "حالة الشيكات"},
		{2, "Exemption Status", "حالة الاعفاء"},
		{3, "Academic Advising Category", "حالات الأرشاد الاكاديمى"},
		{4, "Certificate Recieved Types", "انواع استلام الشهادات"},
		{5, "OnlinePayment Request Sourse", "OnlinePayment Request Sourse"},
		{6, "OnlinePayment Channel Mode", "OnlinePayment Channel Mode"},
		{7, "OnlinePayment Payment Mode", "OnlinePayment Payment Mode"},
		{8, "OnlinePayment CustumerIDType", "OnlinePayment CustumerIDType"},
		{9, "OnlinePayment Request Status", "OnlinePayment Request Status"},
		{10, "OnlinePayment Enquery Response Status", "OnlinePayment Enquery Response Status"},
		{11, "Transcript Reasons", "Transcript Reasons"},
		{12, "Change Fees Rule Actions", "Change Fees Rule Actions"},
		{13, "ADD/DRP ACTIONS", "ADD/DRP ACTIONS"},
		{14, "RECEIPTS ACTIONS", "RECEIPTS ACTIONS"},
		{15, "REFUND ACTIONS", "REFUND ACTIONS"},
		{16, "Locker Request Status", "Locker Request Status"},
		{17, "Car Stiker lost Request Status", "Car Stiker lost Request Status"},
		{18, "APPLICANT APPLIED STATUS", "APPLICANT APPLIED STATUS"},
		{19, "ADMISSION ACCEPTANCE METHOD", "ADMISSION ACCEPTANCE METHOD"},
		{20, "Capacity Plan Status", "Capacity Plan Status"},
		{21, "Nationalities Band", "بند الجنسيات"},
		{22, "Study Countries Band", "بند بلاد التدريس"},
		{23, "CaseStudyRequestStatus", "CaseStudyRequestStatus"},
		{24, "Cancel Admission Request Status", "حالات طلبات الغاء القبول"},
		{25, "Under Load Approval", "قبول تحت الضغط الاكاديمي"},
		{26, "Student Fees Refund Request Status", "Student Fees Refund Request Status"},
		{27, "Admission Document Type", "انواع ملفات القبول"},
		{28, "Admission Request Type", "انواع تقديم القبول"},
		{29, "Refund Reasons", "Refund Reasons"},
		{30, "Staff Load Status", "Staff Load Status"},
		{31, "Staff Load Approval", "Staff Load Approval"},
		{32, "Case Study Types", "Case Study Types"},
		{33, "PG Admission Categories", "انواع قبول الدراسات العليا"},
		{34, "EGA_Verification_Document_Status", "حالة تاكيد الملفات فالحكومة الالكترونية"},
		{44, "Transfer Subtypes", "Transfer Subtypes"},
		{45, "Admission University Category", "تصنيف جامعات القبول"},
		{46, "Study Language", "لغة الدراسة"},
		{47, "Certificate Expiration Type", "نوع انتهاء الافادة"},
		{48, "To Enrollment Configration", "الى حالة قيد"},
		{49, "UG English Exam Status", "حالات امتحان اللغة الانجليزية"},
		{50, "Transcript versions", "انواع نسخ الافادات"},
		{51, "Disability", "Disability"},
	}

	stmt, err := db.Prepare(`
		INSERT OR IGNORE INTO ED_CODE_STATUS_CAT (ED_CODE_STATUS_CAT_ID, DESCR_EN, DESCR_AR)
		VALUES (?, ?, ?);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, r := range records {
		if _, err := stmt.Exec(r.ID, r.DescrEn, r.DescrAr); err != nil {
			return err
		}
	}

	log.Println("✅ ED_CODE_STATUS_CAT table populated successfully")
	return nil
}

func SeedEdCodeStatus() error {
	if db == nil {
		log.Fatal("Database not initialized. Call SetDB first.")
	}

	type Row struct {
		ID     int
		En     string
		Ar     string
		Order  int
		CatID  int
	}

	rows := []Row{
		{1, "created", "تم الانشاء", 1, 1},
		{2, "Ready For Submission", "جاهز للارسال", 2, 1},
		{3, "Recieved", "تم الاستلام", 3, 1},
		{4, "Cleared", "تم الاستلام من البنك", 4, 1},
		{5, "Canceled", "ملغي", 0, 1},
		{6, "Exemption Continued", "استمرار اعغاء", 1, 2},
		{7, "Exemption Temporery Stop", "ايقاف مؤقت للاعفاء", 2, 2},
		{8, "Exemption Permenant Stop", "ايقاف نهائي للاعفاء", 3, 2},
		{9, "At risk", "فى خطر", 1, 3},
		{10, "Normal", "طالب عادي", 2, 3},
		{11, "Need to Define Minor", "يحتاج لإضافة تخصص فرعي", 3, 3},
		{12, "Graduation Ceremony", "حفل التخرج", 1, 4},
		{13, "Student Recieved", "أستلام من الطالب", 2, 4},
		{14, "Responsible Recieved", "أستلام من ينوب عنه", 3, 4},
		{15, "Canceled", "ألغيت", 4, 4},
		{16, "WEB", "WEB", 1, 6},
		{17, "SMS", "SMS", 2, 6},
		{18, "WAP", "WAP", 3, 6},
		{19, "MOBILEAPP", "MOBILEAPP", 4, 6},
		{20, "MOBILEAPP", "MOBILEAPP", 5, 6},
		{21, "KIOSK", "KIOSK", 6, 6},
		{22, "SIS-OTHER", "SIS-OTHER", 1, 5},
		{23, "Purchase", "Purchase", 1, 7},
		{24, "Auth/Capture", "Auth/Capture", 2, 7},
		{25, "CPR", "Personal Number", 1, 8},
		{26, "CR", "Commercial Registration Number", 2, 8},
		{27, "SID", "Student ID", 3, 8},
		{28, "BNRN", "Birth Notification Request Number", 4, 8},
		{29, "CCPR", "Child Personal Number", 5, 8},
		{30, "SN", "Subscriber Number", 6, 8},
		{31, "VN", "Vehicle Number", 7, 8},
		{32, "PN", "Passport Number", 8, 8},
		{33, "UN", "Username", 9, 8},
		{34, "FRID", "Forum Registration ID", 10, 8},
		{35, "GCCID", "GCC Identity Number", 11, 8},
		{36, "New Request", "New Request", 1, 9},
		{38, "IntialPaymentRespose", "IntialPaymentRespose", 2, 9},
		{39, "PaymentRedirectRequest", "PaymentRedirectRequest", 3, 9},
		{40, "PaymentNotificationRequest", "PaymentNotificationRequest", 4, 9},
		{41, "PaymentNotificationResponse", "PaymentNotificationResponse", 5, 9},
		{42, "paymentRedirectResponse", "paymentRedirectResponse", 5, 9},
		{43, "PaymentEnquiryRequest", "PaymentEnquiryRequest", 6, 9},
		{44, "PaymentEnquiryResponse", "PaymentEnquiryResponse", 7, 9},
		{45, "ReceiptConfirmationRequest", "ReceiptConfirmationRequest", 8, 9},
		{46, "ReceiptConfirmationResponse", "ReceiptConfirmationResponse", 9, 9},
		{47, "Pending/Unknown", "Pending/Unknown", 1, 10},
		{48, "Incomplet(Trans NotFound)", "Incomplet(Trans NotFound)", 2, 10},
		{49, "Incomplet(Not Proceed)", "Incomplet(Not Proceed)", 3, 10},
		{50, "Complete Successfully", "Complete Successfully", 4, 10},
		{51, "Faild", "Faild", 5, 10},
		{53, "UOB_PORT_ONLINE_PAYMENT", "UOB_PORT_ONLINE_PAYMENT", 2, 5},
		{54, "Official Entity Outside University", "جهة رسمية خارج الجامعة", 1, 11},
		{55, "Academic Department", "القسم الاكاديمي", 2, 11},
		{56, "Official Entity Inside Bahrain", "جهة رسمية داخل البحرين", 3, 11},
		{57, "Official Entity Outside Bahrain", "جهة رسمية خارج البحرين", 4, 11},
		{58, "Practical Training Department", "قسم التدريب العملى", 5, 11},
		{59, "Add Sponsor", "Add Sponsor", 1, 12},
		{60, "Delete Sponsor", "Delete Sponsor", 2, 12},
		{61, "Change Sponsor", "Change Sponsor", 3, 12},
		{62, "Add Nationality", "Add Nationality", 4, 12},
		{63, "Delete Nationality", "Delete Nationality", 5, 12},
		{64, "Change Nationality", "changeNationality", 6, 12},
		{65, "Change Major", "Change Major", 7, 12},
		{66, "change course Amount", "change course Amount", 8, 13},
		{67, "Add Course", "Add Course", 1, 13},
		{68, "Drop Course", "Drop Course", 2, 13},
		{69, "Delete Course", "Delete Course", 3, 13},
		{70, "Update Course", "Update Course", 4, 13},
		{80, "update Receipt", "update Receipt", 1, 14},
		{81, "Add Receipt", "Add Receipt", 2, 14},
		{82, "Delete Receipt", "Delete Receipt", 3, 14},
		{83, "Cancel Receipt", "Cancel Receipt", 4, 14},
		{84, "Refund", "Refund", 1, 15},
		{85, "Cancel Refund", "Cancel Refund", 2, 15},
		{86, "CHANGE SPO RATIO", "CHANGE SPO RATIO", 1, 12},
		{87, "New", "جديد", 1, 16},
		{88, "Paid", "مدفوع", 2, 16},
		{89, "Ready to Recieve Key", "جاهز لاستلام المفتاح", 3, 16},
		{90, "Recieved Key", "أستلم المفتاح", 4, 16},
		{91, "Returned Key", "رجع المفتاح", 5, 16},
		{93, "NewRequest", "NewRequest", 1, 17},
		{94, "ReadyForPay", "ReadyForPay", 2, 17},
		{95, "Paied", "Paied", 3, 17},
		{96, "ReadyForRecieved", "ReadyForRecieved", 4, 17},
		{97, "Achieved", "Achieved", 5, 17},
		{98, "Refused", "Refused", 6, 17},
		{99, "NEW IGA APPLICANTS", "NEW IGA APPLICANTS", 1, 18},
		{100, "Exceptional APPLICANTS for applied", "Exceptional APPLICANTS for applied", 2, 18},
		{101, "Manual Enroll.", "Manual Enroll.", 1, 19},
		{102, "Auto Enroll.", "Auto Enroll.", 2, 19},
		{103, "NEW", "جديد", 1, 20},
		{104, "Approved By Head of Department", "معتمد من رئيس القسم", 2, 20},
		{105, "Rejected by department head", "طلب مروفوض من رئيس القسم الأكاديمي", 3, 20},
		{106, "Approve By College Dean", "معتمد من عميد الكلية", 4, 20},
		{107, "Rejected by college dean", "طلب مرفوض من عميد الكلية", 5, 20},
		{108, "Approved by admission ", "طلب معتمد من القبول", 6, 20},
		{109, "Rejected by admission ", "طلب مروفوض من القبول", 7, 20},
		{110, "ALL", "الكل", 1, 21},
		{111, "Bahraini", "بحرينى", 2, 21},
		{112, "GCC", "خليجى", 3, 21},
		{113, "Others", "أخرى", 4, 21},
		{114, "ALL", "الكل", 1, 22},
		{115, "Bahraini", "بحرينى", 2, 22},
		{116, "GCC", "خليجى", 3, 22},
		{117, "Others", "أخرى", 4, 22},
		{118, "Pending for Payment", "Pending for Payment", 1, 23},
		{119, "New Request", "New Request", 2, 23},
		{120, "Back To Student", "Back To Student", 3, 23},
		{121, "Approved by HOD", "Approved by HOD", 4, 23},
		{122, "Rejected by HOD", "Rejected by HOD", 5, 23},
		{123, "Approved by Dean of College", "Approved by Dean of College", 6, 23},
		{124, "Approved by Commitee", "Approved by Commitee", 7, 23},
		{125, "Rejected by Commitee", "Rejected by Commitee", 8, 23},
		{126, "Forwarded to Dean of Students Affairs", "Forwarded to Dean of Students Affairs", 9, 23},
		{127, "Approved by Dean of Students Affairs", "Approved by Dean of Students Affairs", 10, 23},
		{128, "Rejected by Dean of Students Affairs", "Rejected by Dean of Students Affairs", 11, 23},
		{129, "Approved by the University Council", "Approved by the University Council", 12, 23},
		{130, "Rejected by the University Council", "Rejected by the University Council", 13, 23},
		{131, "Action Taken by Registeration", "Action Taken by Registeration", 14, 23},
		{132, "Rejected by Registeration", "Rejected by Registeration", 15, 23},
		{133, "New Request", "طلب جديد", 1, 24},
		{134, "Approved by Admission Head", "قبول من رئيس القبول", 2, 24},
		{135, "Rejected by Admission Head", "رفض من رئيس القبول", 3, 24},
		{136, "Send To Dean", "Send To Dean", 16, 23},
		{137, "Arab", "عربى", 5, 21},
		{138, "Foreigner", "اجنبى", 6, 21},
		{139, "Arab", "عربى", 5, 22},
		{140, "Foreigner", "اجنبى", 6, 22},
		{141, "Assigned to Staff", "Assigned to Staff", 17, 23},
		{142, "Approved By Staff", "Approved By Staff", 18, 23},
		{143, "Rejected By Staff", "Rejected By Staff", 19, 23},
		{144, "Approved By Pg Dean", "Approved By Pg Dean", 20, 23},
		{145, "Rejected By Pg Dean", "Rejected By Pg Dean", 21, 23},
		{146, "Rejected By College Dean", "Rejected By College Dean", 22, 23},
		{147, "Send To HOD", "Send To HOD", 22, 23},
		{148, "Ignore current Status", "تجاهل الوضع الحالي", 23, 25},
		{149, "Approved", "Approved", 1, 26},
		{150, "Rejected", "Rejected", 1, 26},
		{151, "NEW", "NEW", 1, 26},
		{152, "Morning Document", "ملفات الصباحى", 1, 27},
		{153, "Evening Document", "ملفات المسائى", 2, 27},
		{154, "Midwifery Document", "ملفات برنامج القبالة ", 3, 27},
		{155, "UG MORNING PROGRAMS", "بكالوريوس الفترة الصباحية", 1, 28},
		{156, "PG MORNING PROGRAMS", "دراسات عليا الفترة الصباحية", 2, 28},
		{157, "UG EVENING PROGRAMS", "بكالوريوس الفترة المسائية", 3, 28},
		{158, "UG MIDWIFERY PROGRAMS", "بكالوريوس برنامج القبالة", 4, 28},
		{159, "UG DIPLOMA TO BSC", "تقديم من الدبلوم الى البكالوريوس", 5, 28},
		{160, "UG COMPLETION PROGRAMS", "بكالوريوس البرامج التكميلية", 6, 28},
		{161, "UG EXTERNAL TRANSFER", "تحويل خارجى", 7, 28},
		{162, "UG_VISITOR_STUDENT", "طالب زائر", 8, 28},
		{163, "Maximum Semester Warning", "تنبيه المدة القصوى", 4, 3},
		{165, "Canceled", "Canceled", 23, 23},
		{166, "ابتعاث", "Sponsorship", 1, 29},
		{167, "إعفاء", "Exemption", 2, 29},
		{168, "تظلم", "Grade Review", 3, 29},
		{169, "حالة دراسية", "Case Study", 4, 29},
		{170, "معاملة كبحريني", "Bahrain similarity", 5, 29},
		{171, "حذف مقرر", "Drop Course", 6, 29},
		{172, "انسحاب", "Withdrawn", 7, 29},
		{173, "تخرج", "Graduation", 8, 29},
		{174, "أخرى", "Others", 9, 29},
		{175, "Updated", "Updated", 3, 26},
		{176, "Completed", "Completed", 4, 26},
		{177, "Normal Load", "Normal Load", 1, 30},
		{178, "UnderLoad", "UnderLoad", 2, 30},
		{179, "OverLoad", "OverLoad", 3, 30},
		{180, "Validation Status", "Validation Status", 1, 31},
		{181, "Approved By HOD", "Approved By HOD", 2, 31},
		{182, "Rejected by Dean", "Rejected by Dean", 3, 31},
		{183, "Approved by Dean", "Approved by Dean", 4, 31},
		{184, "Approved By Council", "Approved By Council", 5, 31},
		{185, "Approved by HR", "Approved by HR", 6, 31},
		{186, "Paid", "Paid", 7, 31},
		{187, "SendBackToDean", "SendBackToDean", 8, 31},
		{188, "Case Study (UG/PG Default)", "Case Study (UG/PG Default)", 1, 32},
		{189, "Case Study(Associate Diploma Dismissed)", "Case Study (Associate Diploma Dismissed)", 2, 32},
		{190, "Case Study(Associate Diploma Enrolled)", "Case Study (Associate Diploma Enrolled)", 3, 32},
		{191, "Case Study (PG Higher Diploma)", "Case Study (PG Higher Diploma)", 4, 32},
		{195, "Undo Hr Approval", "Undo Hr Approval", 8, 31},
		{201, "Conditional Admission on passing the English language", "قبول مشروط باجتياز اللغة الانجليزية", 1, 33},
		{202, "Conditional Admission on passing specialized remedial courses and passing the English language", "قبول مشروط باجتياز المقررات الاستدراكية التخصصية والمجتاز لشرط اللغة الانجليزية", 2, 33},
		{203, "Conditional Admission with less than the required GPA and passing the English language requirement", "قبول مشروط بأقل من المعدل المطلوب والمجتاز لشرط اللغة الانجليزية", 3, 33},
		{204, "Conditional Admission with less than the required GPA and did not pass the English language requirement", "قبول مشروط بأقل من المعدل المطلوب ولم يجتاز لشرط اللغة الانجليزية", 4, 33},
		{205, "Conditional Admission on passing the English language + specialized remedial courses", "قبول مشروط باجتياز اللغة الانجليزية والمواد الاستدراكية التخصصية", 5, 33},
		{206, "Unconditional Admission", "قبول غير مشروط", 6, 33},
		{207, "Default", "Default", 1, 44},
		{208, "Applied 2023", "Applied 2023", 2, 44},
		{209, "Generate Qr Code", "Generate Qr Code", 1, 34},
		{210, "Activated", "Activated", 2, 34},
		{211, "Canceled", "Canceled", 3, 34},
		{212, "Suspended", "Suspended", 4, 34},
		{213, "Graduate from the University of Bahrain", "خريج من جامعة البحرين", 1, 45},
		{214, "Graduate from other universities in Bahrain", "خريج من جامعات أخرى في البحرين", 2, 45},
		{215, "Graduate from universities outside Bahrain", "خريج من جامعات خارج البحرين", 3, 45},
		{216, "English", "اللغة الإنجليزية", 1, 46},
		{217, "Arabic", "اللغة العربية", 2, 46},
		{218, "Other", "أخرى", 3, 46},
		{219, "Approved by Student", "قبول من الطالب", 4, 24},
		{220, "by days", "بالايام", 1, 47},
		{221, "by end of current semester", "بانتهاء الفصل الدراسى الحالى", 2, 47},
		{222, "Enrolled", "مقيد", 1, 48},
		{223, "OFFICIAL WITHDRAWAL", "انسحاب رسمى", 2, 48},
		{224, "UNOFFICIAL WITHDRAWAL", "انسحاب غير  رسمى", 3, 48},
		{225, "PERMANENTLY WITHDRAWN FROM UNIVERSITY", "انسحاب نهائى", 4, 48},
		{226, "SUSPENDED", "ايقاف عن الدراسة لمدة فصل اكاديمي واحد", 5, 48},
		{227, "DISMISSED", "فصل من الجامعة ", 6, 48},
		{228, "GRADUATED", "تخرج", 7, 48},
		{229, "MAXIMUM SEMESTERS DISMISSAL", "فصل نهائي بسبب انتهاء المدة القصوى", 8, 48},
		{230, "DISMISSED FROM BTC COLLEGE", "مفصول من كلية البحرين للمعلمين", 9, 48},
		{231, "CANCEL REGISTRATION", "إلغاء تسجيل", 10, 48},
		{232, "OFFICIAL WITHDRAWAL (TWO SEMESTERS)", "انسحاب رسمى لمدة فصلين", 11, 48},
		{233, "Fulfilled Graduation Requirements - Pending Council Approval", "تم استيفاء متطلبات التخرج وفي انتظار اعتماد مجلس الجامعة", 12, 48},
		{234, "DISCIPLINARY DISMISSAL", "فصل تأديبي نهائي", 13, 48},
		{235, "DISCIPLINARY SUSPENDED 2 SEMESTERS", "فصل تأديبي لمدة فصلين", 14, 48},
		{236, "DISCIPLINARY SUSPENSION (ONE SEMESTER)", "فصل تأديبي لفصل واحد", 15, 48},
		{237, "Cancel Admission", "الغاء قبول", 16, 48},
		{238, "Cancel ReEnroll", "الغاء قيد", 17, 48},
		{239, "All courses arenot active", "كل المقررات غير نشطة", 18, 48},
		{240, "Conditional Admission on Passing Specialized Remedial Courses", "قبول مشروط باجتياز المقررات الاستدراكية التخصصية", 7, 33},
		{241, "Executed", "تم التنفيذ", 1, 49},
		{242, "Not Executed", "لم يتم التنفيذ", 2, 49},
		{243, "Added to Student Balance", "تم الاضافة لرصيد الطالب الجامعي", 5, 26},
		{244, "Applied >=2024", "Applied >=2024", 3, 44},
		{245, "Electronic", "الكترونية ", 1, 50},
		{246, "Printed", "مطبوعة", 2, 50},
		{247, "New", "New", 1, 51},
		{248, "Created", "Created", 2, 51},
		{249, "Approved", "Approved", 3, 51},
		{250, "Rejected", "Rejected", 4, 51},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO ED_CODE_STATUS
			(ED_CODE_STATUS_ID, DESCR_EN, DESCR_AR, STATUS_ORDER, ED_CODE_STATUS_CAT_ID)
		VALUES (?, ?, ?, ?, ?);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, r := range rows {
		if _, err := stmt.Exec(r.ID, r.En, r.Ar, r.Order, r.CatID); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Println("✅ ED_CODE_STATUS table populated successfully")
	return nil
}

func SeedGsCodeReqStatus() error {
	if db == nil {
		log.Fatal("Database not initialized. Call SetDB first.")
	}

	type Row struct {
		ID     int
		Ar     string
		En     string
	}

	rows := []Row{
		{1, "طلب جديد", "New"},
		{2, "موافقة", "Approved"},
		{3, "مرفوض", "Refused"},
		{4, "رفض من مدرس المادة", "Rejected by Staff Member"},
		{5, "رفض من مدير القسم", "Rejected by Head of Dep."},
		{6, "رفض من مديرالتسجيل", "Rejected by Registration"},
		{7, "الغاء", "Cancel"},
		{9, "قبول من المرشد الاكاديمي", "Approved by Advisor"},
		{10, "قبول من مدير التسجيل", "Approved by Registration"},
		{11, "قبول من مدير القسم", "Approved by Head of Dep"},
		{12, "رفض من المرشد الاكاديمي", "Rejected by Advisor"},
		{13, "رفض بعد غلق السمستر", "Rejected after closing"},
		{14, "Approved by Force", "Approved by Force"},
		{15, "Not Complete", "Not Complete"},
		{16, "تم الإعتماد من اللجنة", "Approved by committe"},
		{17, "تم الاعتماد من العميد", "Approved by dean"},
		{18, "تم إدخال الدرجة", "Grade Entered"},
		{19, "رفض من المشرف", "Rejected by Supervisor"},
		{20, "قبول من المرشد", "Approved by Supervisor"},
		{21, "Under Supervisor", "Under Supervisor"},
		{22, "Approved By Reviewer", "Approved By Reviewer"},
		{23, "Rejected By Reviewer", "Rejected By Reviewer"},
		{24, "Under Reviewer", "Under Reviewer"},
		{25, "Under Head Of Department", "Under Head Of Department"},
		{26, "رفض من النظام", "Rejected by System"},
		{27, "REJECT_WITH_REFUND", "REJECT_WITH_REFUND"},
		{28, "تم تغيير الدرجة للطالب", "Grade changed for student"},
		{29, "ClosedByRegistration", "ClosedByRegistration"},
		{30, "OpendByRegistration", "OpendByRegistration"},
		{31, "Under Referrer", "Under Referrer"},
		{32, "Approved By Referrer", "Approved By Referrer"},
		{33, "Return To Reviewers", "Return To Reviewers"},
		{34, "في انتظار المستند", "Pending for Document"},
		{35, "Approve No Change", "Approve No Change"},
		{36, "Approve With Change", "Approve With Change"},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO GS_CODE_REQ_STATUS
			(GS_CODE_REQ_STATUS_ID, DESCR_AR, DESCR_EN)
		VALUES (?, ?, ?);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, r := range rows {
		if _, err := stmt.Exec(r.ID, r.Ar, r.En); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Println("✅ GS_CODE_REQ_STATUS table populated successfully")
	return nil
}



