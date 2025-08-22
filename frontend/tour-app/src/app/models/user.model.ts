export interface User {
  id: number; // Ovo polje zavisi od toga šta vaš backend vraća u JWT tokenu
  username: string;
  email: string;
  role: 'turista' | 'vodic' | 'administrator';
  isBlocked: boolean;

  // Opciona polja za profil (koristimo '?' da označimo da ne moraju uvek postojati)
  firstName?: string;
  lastName?: string;
  profilePicture?: string;
  biography?: string;
  motto?: string;
}