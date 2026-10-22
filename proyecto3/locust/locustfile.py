from locust import HttpUser, task, between
import random
import json
from datetime import datetime
#paises permitidos 

COUNTRIES = ["USA", "RUS", "CHN", "ESP", "GTM"]

class MilitaryReporter(HttpUser):
    #tiempo de espera entre peticiones (entre 0.5 y 1.5 seg)
    wait_time = between(0.5, 1.5)


    @task
    #definimos el reporte 
    def sedReport(self):
        #generamos un reporte aleatorio
        report = {
            "country": random.choice(COUNTRIES),
            "warplanes_in_air": random.randint(0, 50),
            "warships_in_water": random.randint(0, 30),
            "timestamp": datetime.now().isoformat()

        }

        #esto se lo mandamos al endpoint /report
        with self.client.post(
            "/report",
            json=report,
            catch_response=True
        ) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"Failed to send report: {response.status_code} - {response.text}")
                