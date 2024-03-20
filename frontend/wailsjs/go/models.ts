export namespace canvas {
	
	export enum CourseEnrollmentType {
	    TEACHER = "teacher",
	    STUDENT = "student",
	    TA = "ta",
	    OBSERVER = "observer",
	    DESIGNER = "designer",
	}
	export enum EnrollmentType {
	    TEACHER = "TeacherEnrollment",
	    STUDENT = "StudentEnrollment",
	    TA = "TaEnrollment",
	    DESIGNER = "DesignerEnrollment",
	    OBSERVER = "ObserverEnrollment",
	}
	export enum AssignmentBucket {
	    PAST = "past",
	    OVERDUE = "overdue",
	    UPDATED = "undated",
	    UNGRADED = "ungraded",
	    UNSUBMITTED = "unsubmitted",
	    UPCOMING = "upcoming",
	    FUTURE = "future",
	}
	export class Account {
	    id: number;
	    name: string;
	    parent_account_id: number;
	    root_account_id: number;
	    workflow_state: string;
	
	    static createFrom(source: any = {}) {
	        return new Account(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.parent_account_id = source["parent_account_id"];
	        this.root_account_id = source["root_account_id"];
	        this.workflow_state = source["workflow_state"];
	    }
	}
	export class AssignmentDate {
	    id: number;
	    due_at: string;
	    unlock_at: string;
	    lock_at: string;
	    title: string;
	    set_type: string;
	    set_id: number;
	
	    static createFrom(source: any = {}) {
	        return new AssignmentDate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.due_at = source["due_at"];
	        this.unlock_at = source["unlock_at"];
	        this.lock_at = source["lock_at"];
	        this.title = source["title"];
	        this.set_type = source["set_type"];
	        this.set_id = source["set_id"];
	    }
	}
	export class SectionNeedsGrading {
	    section_id: number;
	    needs_grading_count: number;
	
	    static createFrom(source: any = {}) {
	        return new SectionNeedsGrading(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.section_id = source["section_id"];
	        this.needs_grading_count = source["needs_grading_count"];
	    }
	}
	export class Assignment {
	    id: number;
	    course_id: number;
	    qualification: string;
	    course_name: string;
	    name: string;
	    due_at: string;
	    unlock_at: string;
	    lock_at: string;
	    needs_grading_count: number;
	    section: string;
	    needs_grading_section: number;
	    teachers: string;
	    status: string;
	    published: boolean;
	    gradebook_url: string;
	    needs_grading_count_by_section: SectionNeedsGrading[];
	    all_dates: AssignmentDate[];
	
	    static createFrom(source: any = {}) {
	        return new Assignment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.course_id = source["course_id"];
	        this.qualification = source["qualification"];
	        this.course_name = source["course_name"];
	        this.name = source["name"];
	        this.due_at = source["due_at"];
	        this.unlock_at = source["unlock_at"];
	        this.lock_at = source["lock_at"];
	        this.needs_grading_count = source["needs_grading_count"];
	        this.section = source["section"];
	        this.needs_grading_section = source["needs_grading_section"];
	        this.teachers = source["teachers"];
	        this.status = source["status"];
	        this.published = source["published"];
	        this.gradebook_url = source["gradebook_url"];
	        this.needs_grading_count_by_section = this.convertValues(source["needs_grading_count_by_section"], SectionNeedsGrading);
	        this.all_dates = this.convertValues(source["all_dates"], AssignmentDate);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AssignmentResult {
	    assignment_id: number;
	    user_sis_id: string;
	    stduent_name: string;
	    qualification: string;
	    course_name: string;
	    title: string;
	    max_score: number;
	    min_score: number;
	    // Go type: struct { Score float32 "json:\"score\" csv:\"Score\""; SubmittedAt string "json:\"submitted_at\" csv:\"Submitted At\"" }
	    submission: any;
	    "due-at": string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new AssignmentResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.assignment_id = source["assignment_id"];
	        this.user_sis_id = source["user_sis_id"];
	        this.stduent_name = source["stduent_name"];
	        this.qualification = source["qualification"];
	        this.course_name = source["course_name"];
	        this.title = source["title"];
	        this.max_score = source["max_score"];
	        this.min_score = source["min_score"];
	        this.submission = this.convertValues(source["submission"], Object);
	        this["due-at"] = source["due-at"];
	        this.status = source["status"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Course {
	    id: number;
	    name: string;
	    course_code: string;
	    account_id: number;
	    root_account_id: number;
	    friendly_name: string;
	    workflow_state: string;
	    start_at: string;
	    end_at: string;
	    is_public: boolean;
	    enrollment_term_id: number;
	    account?: Account;
	
	    static createFrom(source: any = {}) {
	        return new Course(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.course_code = source["course_code"];
	        this.account_id = source["account_id"];
	        this.root_account_id = source["root_account_id"];
	        this.friendly_name = source["friendly_name"];
	        this.workflow_state = source["workflow_state"];
	        this.start_at = source["start_at"];
	        this.end_at = source["end_at"];
	        this.is_public = source["is_public"];
	        this.enrollment_term_id = source["enrollment_term_id"];
	        this.account = this.convertValues(source["account"], Account);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Enrollment {
	    id: number;
	    user_id: number;
	    course_id: number;
	    course_section_id: number;
	    sis_section_id: string;
	    // Go type: struct { HtmlUrl string "json:\"html_url\""; CurrentScore float32 "json:\"current_score\""; CurrentGrade string "json:\"current_grade\""; FinalScore float32 "json:\"final_score\""; FinalGrade string "json:\"final_grade\"" }
	    grades: any;
	    // Go type: struct { Name string "json:\"name\""; SISUserID string "json:\"sis_user_id\"" }
	    user: any;
	    enrollment_state: string;
	    role: string;
	
	    static createFrom(source: any = {}) {
	        return new Enrollment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.user_id = source["user_id"];
	        this.course_id = source["course_id"];
	        this.course_section_id = source["course_section_id"];
	        this.sis_section_id = source["sis_section_id"];
	        this.grades = this.convertValues(source["grades"], Object);
	        this.user = this.convertValues(source["user"], Object);
	        this.enrollment_state = source["enrollment_state"];
	        this.role = source["role"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EnrollmentResult {
	    sis_id: string;
	    name: string;
	    account: string;
	    course_name: string;
	    section: string;
	    enrollment_state: string;
	    course_state: string;
	    current_grade: string;
	    current_score: number;
	    enrollment_role: string;
	    gardes_url: string;
	
	    static createFrom(source: any = {}) {
	        return new EnrollmentResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sis_id = source["sis_id"];
	        this.name = source["name"];
	        this.account = source["account"];
	        this.course_name = source["course_name"];
	        this.section = source["section"];
	        this.enrollment_state = source["enrollment_state"];
	        this.course_state = source["course_state"];
	        this.current_grade = source["current_grade"];
	        this.current_score = source["current_score"];
	        this.enrollment_role = source["enrollment_role"];
	        this.gardes_url = source["gardes_url"];
	    }
	}
	
	export class Section {
	    id: number;
	    sis_section_id: string;
	    name: string;
	    start_at: string;
	    end_at: string;
	    course_id: number;
	    total_students: number;
	
	    static createFrom(source: any = {}) {
	        return new Section(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.sis_section_id = source["sis_section_id"];
	        this.name = source["name"];
	        this.start_at = source["start_at"];
	        this.end_at = source["end_at"];
	        this.course_id = source["course_id"];
	        this.total_students = source["total_students"];
	    }
	}
	export class Submission {
	    id: number;
	    // Go type: struct { SISUserID string "json:\"sis_user_id\" csv:\"ID\""; Name string "json:\"name\" csv:\"Name\"" }
	    user: any;
	    user_id: number;
	    assignment_id: number;
	    assignment_name: string;
	    course_id: number;
	    grade: string;
	    submitted_at: string;
	    graded_at: string;
	    attempt: number;
	    grader_id: number;
	    late: boolean;
	    excused: boolean;
	    preview_url: string;
	
	    static createFrom(source: any = {}) {
	        return new Submission(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.user = this.convertValues(source["user"], Object);
	        this.user_id = source["user_id"];
	        this.assignment_id = source["assignment_id"];
	        this.assignment_name = source["assignment_name"];
	        this.course_id = source["course_id"];
	        this.grade = source["grade"];
	        this.submitted_at = source["submitted_at"];
	        this.graded_at = source["graded_at"];
	        this.attempt = source["attempt"];
	        this.grader_id = source["grader_id"];
	        this.late = source["late"];
	        this.excused = source["excused"];
	        this.preview_url = source["preview_url"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class User {
	    id: number;
	    name: string;
	    sis_user_id: string;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.sis_user_id = source["sis_user_id"];
	    }
	}

}

