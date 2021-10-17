# MyCal

MyCal est un outil pour crée un fichier .ICS (iCalendar) de son emploi du temps MyGes, pour ensuite l'importer dans Google Calendar ou n'importe quel calendrier.

# Utilisation

Il suffit d'aller dans les [release](https://github.com/obito/mycal/releases) et de télécharger l'éxécutable pour votre plateforme. Ensuite, connectez-vous avec vos identifiants **MyGES**

Le format utilisé pour les dates est mm/dd/yyyy (mois/jour/année, format US).

# Comment?

Il est possible d'utiliser cet outil même lorsque MyGES n'est pas fonctionel. J'ai découvert que l'application Skolae avait leur propre endpoints sur l'api de Kordis (ce qu'utilise MyGES), et donc ne sont donc pas relié à MyGES. Pour comprendre et utiliser ces endpoints, j'ai simplement mis en place un proxy grâce à Burp Suite, puis configurer mon iPhone pour utiliser ce proxy. 
