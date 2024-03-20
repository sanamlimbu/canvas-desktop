// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {canvas} from '../models';

export function GetAccountByID(arg1:number):Promise<canvas.Account>;

export function GetAssignmentsByAccount(arg1:canvas.Account,arg2:canvas.AssignmentBucket):Promise<Array<canvas.Assignment>>;

export function GetAssignmentsByCourse(arg1:canvas.Course,arg2:canvas.AssignmentBucket):Promise<Array<canvas.Assignment>>;

export function GetAssignmentsResultsByUser(arg1:canvas.User):Promise<Array<canvas.AssignmentResult>>;

export function GetCourseByID(arg1:number):Promise<canvas.Course>;

export function GetCoursesByAccount(arg1:canvas.Account,arg2:canvas.CourseEnrollmentType):Promise<Array<canvas.Course>>;

export function GetCoursesByAccountID(arg1:number,arg2:canvas.CourseEnrollmentType):Promise<Array<canvas.Course>>;

export function GetCoursesByUser(arg1:canvas.User):Promise<Array<canvas.Course>>;

export function GetEnrollmentResultsByUser(arg1:canvas.User):Promise<Array<canvas.EnrollmentResult>>;

export function GetEnrollmentsBySectionID(arg1:number,arg2:Array<canvas.EnrollmentType>):Promise<Array<canvas.Enrollment>>;

export function GetEnrollmentsByUser(arg1:canvas.User):Promise<Array<canvas.Enrollment>>;

export function GetSectionByID(arg1:number):Promise<canvas.Section>;

export function GetSectionsByCourseID(arg1:number):Promise<Array<canvas.Section>>;

export function GetSubmissions(arg1:number,arg2:number):Promise<Array<canvas.Submission>>;

export function GetUngradedSubmissionsByAccount(arg1:canvas.Account):Promise<Array<canvas.Submission>>;

export function GetUserBySisID(arg1:string):Promise<canvas.User>;
