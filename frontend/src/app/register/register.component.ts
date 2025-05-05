import { Component } from '@angular/core';
import { AuthService } from '../auth.service';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-register',
  imports: [FormsModule, CommonModule, RouterModule],
  standalone: true,
  templateUrl: './register.component.html',
  styleUrl: './register.component.css'
})
export class RegisterComponent {
  username: string = '';
  name: string = '';
  surname: string = '';
  email: string = '';
  password: string = '';
  role: string = 'student'; // Default role
  errorMessage: string = '';

  constructor(private authService: AuthService, private router: Router) {}

  onSubmit() {
    this.errorMessage = '';
    this.authService.register(this.username, this.name, this.surname ,this.email, this.password, this.role).subscribe({
      next: (response) => {
        console.log('Registration successful:', response);
        this.router.navigate(['/login']);
      },
      error: (err) => {
        if (err.status === 400) {
          this.errorMessage = 'Please enter a valid email and password (minimum 6 characters)';
          console.log('Invalid email or password:', err, this.email, this.password, this.username, this.role);
        } else if (err.status === 409) {
          this.errorMessage = 'Email already exists';
        } else {
          this.errorMessage = 'An error occurred. Please try again later.';
        }
        console.error('Registration error:', err);
      },
    });
  }
}
