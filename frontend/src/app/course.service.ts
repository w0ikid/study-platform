import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map, Observable } from 'rxjs';


export interface Course {
  id: number;
  name: string;
  description: string;
  image_url: string;
  teacher_id: number;
  status: 'active' | 'inactive';
  created_at: string;
  updated_at: string;
}

export interface Lesson {
  id: number;
  courseId: number;
  title: string;
  content: string;
  video_url: string;
  created_at: string;
  updated_at: string;
}

export interface EnrollmentResponse {
  message: string;
}

@Injectable({
  providedIn: 'root'
})

export class CourseService {
  private apiUrl = `http://localhost:8080/api/courses/`
  constructor( private http: HttpClient ) { }

  getAllCourses(): Observable<Course[]> {
    return this.http.get<{ courses: Course[] }>(this.apiUrl).pipe(
      map(response => response.courses)
    );
  }
  enroll(courseId: number): Observable<EnrollmentResponse> {
    return this.http.post<EnrollmentResponse>(`${this.apiUrl}${courseId}/enroll/`, {}, { withCredentials: true});
  }

  getCourseById(id: number): Observable<Course> {
    return this.http.get<Course>(`${this.apiUrl}${id}/`, {withCredentials: true});
  }
}
