from urllib import urlopen
from flask import Flask
import re

app = Flask(__name__)

def GetCourseListing(course):
    page = urlopen("http://www.registrar.ufl.edu/cdesc?crs=" + course.upper())

    document = ""

    for line in page.readlines():
        document = document + line + '\n'

    return document

def Parse(document):
    classNo = re.compile(r"<h2>(?P<class>.*)</h2>")
    classTitle = re.compile(r"<h3>(?P<title>.*)</h3>")
    classReq = re.compile(r"<strong>Credits: (?P<credits>\d*); Prereq: (?P<prereq>.*)\.</strong>")

    result = classNo.search(document)
    info = {}

    info["class"] = result.group("class")

    result = classTitle.search(document)
    info["title"] = result.group("title")

    result = classReq.search(document)
    info["credits"] = result.group("credits")
    info["prereq"] = result.group("prereq")

    return info

@app.route("/")
def index():
    page = open("templates/index.html")

    data = page.read()

    page.close()

    return data

def main():
    cop4600 = GetCourseListing("CEG4011")

    cop4600 = Parse(cop4600)

    print cop4600

main()
