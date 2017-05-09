echo off
kubectl.exe --namespace amdr set image deploy/amdr-backend amdr-backend=us.gcr.io/personal-projects-1369/amdr/backend:%1