import { Component, OnInit } from '@angular/core';
import { Course, CourseService } from '../course.service';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';

@Component({
  selector: 'app-course-list',
  standalone: true,
  imports: [CommonModule, MatButtonModule, MatCardModule],
  templateUrl: './course-list.component.html',
  styleUrl: './course-list.component.css'
})
export class CourseListComponent implements OnInit {
  courses: Course[] = [];

  constructor(private courseService: CourseService, private router: Router) {}

  ngOnInit(): void {
    this.loadCourses();
  }
  loadCourses(): void {
    this.courseService.getAllCourses().subscribe({
      next: (courses) => {
        this.courses = courses;
      },
      error: (error) => {
        console.error('Ошибка загрузки курсов:', error);
      },
    });
  }

  enroll(courseId: number): void {
    this.courseService.enroll(courseId).subscribe({
      next: (response) => {
        alert(response.message); // Показываем сообщение от сервера
      },
      error: (error) => {
        console.error('Ошибка при записи на курс:', error);
        alert('Не удалось записаться на курс');
      },
    });
  }

  viewDetails(courseId: number): void {
    this.router.navigate(['/courses', courseId]); // Перенаправление на страницу курса
  }
}
