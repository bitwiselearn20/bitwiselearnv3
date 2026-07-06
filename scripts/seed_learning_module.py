#!/usr/bin/env python3
"""
Build a COMPLETE BitwiseLearn learning module end-to-end through the API gateway.

What it creates (self-contained, runs cleanly every time with unique suffixes):
  - A DSA problem ("Echo Sum") with an example test case + Python template
  - A course + one module/section
  - 5 trackable learning-content items:
      1. YouTube video lecture        (video_url)
      2. PDF reading                  (uploaded file  -> Border_Wall_Estimate.pdf)
      3. Lecture notes                (text/transcript)
      4. Worked-example video         (video_url)
      5. Practice problem             (links to the problem bank item)
  - A 10-question MCQ/SCQ assignment
  - An institution + batch + student, enrolled in the course
Then it acts AS THE STUDENT to prove the completed / not-completed tracking:
  - marks 2 of the 5 items done
  - submits the assignment
  - reads back per-item progress (completed: true/false) and assignment marks

Run:
  uv run --no-project --with requests python scripts/seed_learning_module.py
"""

import sys
import time
import requests

GATEWAY = "http://localhost:8000"
ADMIN_EMAIL = "superadmin@bitwiselearn.com"
ADMIN_PASSWORD = "Admin@123"

YOUTUBE_EMBED = "https://www.youtube.com/embed/4hcRaRultRM?si=FZDvrwndJs-rhqVx"
PDF_PATH = "/Users/brittoanand/Downloads/Border_Wall_Estimate.pdf"

SUFFIX = str(int(time.time()))  # unique per run


def die(msg, resp=None):
    print(f"\nFAILED: {msg}")
    if resp is not None:
        print(f"  HTTP {resp.status_code}: {resp.text[:400]}")
    sys.exit(1)


class API:
    def __init__(self, base):
        self.base = base
        self.token = None

    def _headers(self, extra=None):
        h = {}
        if self.token:
            h["Authorization"] = f"Bearer {self.token}"
        if extra:
            h.update(extra)
        return h

    def post(self, path, json=None, files=None, data=None):
        h = self._headers(None if files else {"Content-Type": "application/json"})
        return requests.post(self.base + path, json=json, files=files, data=data, headers=h, timeout=60)

    def get(self, path):
        return requests.get(self.base + path, headers=self._headers(), timeout=60)

    def put(self, path, json=None):
        return requests.put(self.base + path, json=json, headers=self._headers({"Content-Type": "application/json"}), timeout=60)


def data_of(resp, what):
    try:
        body = resp.json()
    except Exception:
        die(f"{what}: non-JSON response", resp)
    if resp.status_code >= 400 or body.get("statusCode", 200) >= 400:
        die(f"{what}: {body.get('message')}", resp)
    return body.get("data") or {}


def main():
    admin = API(GATEWAY)

    # ---- 0. Admin login -------------------------------------------------
    r = admin.post("/api/v1/auth/admin/login", json={"email": ADMIN_EMAIL, "password": ADMIN_PASSWORD})
    admin.token = (data_of(r, "admin login").get("tokens") or {}).get("accessToken")
    if not admin.token:
        die("no admin token", r)
    print("admin logged in")

    # ---- 1. Problem (for 'solve problem from problem bank') -------------
    pid = data_of(admin.post("/api/v1/problems/add-problem/", json={
        "name": f"Echo Sum {SUFFIX}",
        "description": "Read two integers on one line and print their sum.",
        "difficulty": "EASY",
        "hints": ["Use input().split()"],
    }), "add problem")["id"]
    admin.post(f"/api/v1/problems/add-testcase-to-problem/{pid}",
               json={"testType": "EXAMPLE", "input": "2 3", "output": "5"})
    admin.post(f"/api/v1/problems/add-template-to-problem/{pid}",
               json={"language": "PYTHON",
                     "defaultCode": "a,b=map(int,input().split())\nprint(a+b)",
                     "functionBody": ""})
    print(f"problem created: {pid}")

    # ---- 2. Course + module/section -------------------------------------
    cid = data_of(admin.post("/api/v1/courses/create-course", json={
        "name": f"Python & DSA Foundations {SUFFIX}",
        "description": "A complete starter module: video, reading, notes, practice and a quiz.",
        "level": "BASIC",
        "instructorName": "Dr. Ada Lovelace",
    }), "create course")["id"]
    sid = data_of(admin.post(f"/api/v1/courses/add-course-section/{cid}", json={
        "name": "Module 1: Foundations",
        "description": "Watch, read, practise and test yourself.",
    }), "add section")["id"]
    print(f"course {cid} / section {sid}")

    # ---- 3. Five trackable content items --------------------------------
    def add_content(name, description, video_url="", transcript=""):
        return data_of(admin.post("/api/v1/courses/add-content-to-section", json={
            "sectionId": sid, "name": name, "description": description,
            "videoUrl": video_url, "transcript": transcript,
        }), f"add content '{name}'")["id"]

    c1 = add_content("1. Lecture: Course Introduction (Video)",
                     "Kick-off lecture.", video_url=YOUTUBE_EMBED)
    c2 = add_content("2. Reading: Reference Document (PDF)",
                     "Download and read the reference PDF.")
    c3 = add_content("3. Lecture Notes (Text)",
                     "Key takeaways from the lecture.",
                     transcript="Arrays, time complexity (Big-O), and reading stdin in Python.")
    c4 = add_content("4. Worked Example (Video)",
                     "Walkthrough of a sample problem.", video_url=YOUTUBE_EMBED)
    c5 = add_content("5. Practice: Solve a Problem",
                     f"Open the problem bank and solve it: /problems/{pid}",
                     transcript=f"PROBLEM_ID={pid}")
    contents = [c1, c2, c3, c4, c5]
    print(f"5 content items created: {contents}")

    # Upload the real PDF into content item #2
    try:
        with open(PDF_PATH, "rb") as fh:
            up = admin.post(f"/api/v1/courses/upload-file-in-content/{c2}",
                            files={"file": ("Border_Wall_Estimate.pdf", fh, "application/pdf")})
        file_url = data_of(up, "upload pdf").get("file") or data_of(up, "upload pdf")
        print(f"PDF uploaded to content #2 -> {file_url}")
    except FileNotFoundError:
        print(f"WARN: PDF not found at {PDF_PATH}; skipped upload")

    # ---- 4. Assignment with 10 MCQ/SCQ questions ------------------------
    aid = data_of(admin.post("/api/v1/courses/add-assignment-to-section/", json={
        "sectionId": sid, "name": "Module 1 Quiz",
        "description": "10-question check-your-understanding quiz.",
        "instruction": "Choose the best answer. One attempt allowed.",
        "marksPerQuestion": 1,
    }), "add assignment")["id"]

    questions = [
        ("Big-O of accessing an array element by index?", ["O(1)", "O(n)", "O(log n)", "O(n^2)"], ["O(1)"], "SCQ"),
        ("Which reads a full line from stdin in Python?", ["input()", "print()", "len()", "range()"], ["input()"], "SCQ"),
        ("Big-O of linear search in an unsorted list?", ["O(1)", "O(n)", "O(log n)", "O(n log n)"], ["O(n)"], "SCQ"),
        ("Which are linear data structures? (choose all)", ["Array", "Stack", "Queue", "Graph"], ["Array", "Stack", "Queue"], "MCQ"),
        ("Binary search requires the data to be?", ["Sorted", "Unsorted", "Hashed", "Linked"], ["Sorted"], "SCQ"),
        ("LIFO order is provided by a?", ["Queue", "Stack", "Array", "Tree"], ["Stack"], "SCQ"),
        ("Worst-case Big-O of bubble sort?", ["O(n)", "O(n log n)", "O(n^2)", "O(1)"], ["O(n^2)"], "SCQ"),
        ("Which convert a string to int in Python? (choose all)", ["int(x)", "str(x)", "float(x)", "eval(x)"], ["int(x)", "eval(x)"], "MCQ"),
        ("A hash table average lookup is?", ["O(1)", "O(n)", "O(log n)", "O(n^2)"], ["O(1)"], "SCQ"),
        ("Which structure is FIFO?", ["Stack", "Queue", "Array", "Heap"], ["Queue"], "SCQ"),
    ]
    qids = []
    for q, opts, correct, qtype in questions:
        qid = data_of(admin.post(f"/api/v1/courses/add-assignment-question/{aid}", json={
            "question": q, "options": opts, "correctAnswer": correct, "type": qtype,
        }), "add question")
        qids.append((qid.get("id") if isinstance(qid, dict) else None, correct))
    print(f"assignment {aid} created with {len(questions)} questions")

    # ---- 5. Publish the course so students can see it -------------------
    admin.put(f"/api/v1/courses/change-publish-status/{cid}")
    print("course published")

    # ---- 6. Institution + batch + student, enrolled ---------------------
    iid = data_of(admin.post("/api/v1/institutions/create-institution", json={
        "name": f"Demo College {SUFFIX}", "address": "1 Campus Rd", "pinCode": "560001",
        "email": f"college{SUFFIX}@example.com", "phoneNumber": "9000000000",
    }), "create institution")["id"]
    bid = data_of(admin.post("/api/v1/batches/create-batch", json={
        "batchname": f"CSE-{SUFFIX}", "branch": "CSE", "batchEndYear": "2027", "institutionId": iid,
    }), "create batch")["id"]
    student_email = f"student{SUFFIX}@example.com"
    student_pwd = "Student@123"
    data_of(admin.post("/api/v1/students/create-student", json={
        "name": "Sam Student", "rollNumber": f"R{SUFFIX}", "email": student_email,
        "loginPassword": student_pwd, "batchId": bid, "institutionId": iid,
    }), "create student")
    admin.post("/api/v1/courses/add-course-enrollment/", json={
        "courseId": cid, "batchId": bid, "institutionId": iid,
    })
    print(f"student {student_email} created and batch enrolled")

    # ---- 7. Act as the student ------------------------------------------
    student = API(GATEWAY)
    r = student.post("/api/v1/auth/student/login", json={"email": student_email, "password": student_pwd})
    student.token = (data_of(r, "student login").get("tokens") or {}).get("accessToken")
    if not student.token:
        die("no student token", r)
    print("student logged in")

    # Mark items 1 and 2 as done (2 of 5)
    for c in (c1, c2):
        student.post(f"/api/v1/courses/mark-content-as-done/{c}")
    print("student marked content #1 and #2 as done")

    # Submit the assignment (answer all 10 correctly)
    answers = [{"question_id": qid, "answer": correct} for qid, correct in qids if qid]
    sub = student.post(f"/api/v1/courses/submit-course-assignment/{aid}", json={"answers": answers})
    sub_report = data_of(sub, "submit assignment").get("report", {})

    # ---- 8. Read back the completed / not-completed view ----------------
    # data is a list of sections, each with `contents` -> [{name, completed}].
    sections = data_of(student.get(f"/api/v1/courses/get-individual-course-progress/{cid}"), "progress")

    print("\n================= STUDENT PROGRESS (completed / not-completed) =================")
    for sec in sections:
        print(f"{sec.get('sectionName')}: {sec.get('completed')}/{sec.get('total')} "
              f"items completed ({sec.get('percentage')}%)")
        for item in sec.get("contents", []):
            mark = "[x] DONE    " if item.get("completed") else "[ ] not done"
            print(f"  {mark} {item.get('name')}")
    # report keys are camelCased by the response envelope
    print("\n----------------- ASSIGNMENT RESULT -----------------")
    print(f"  {sub_report.get('assignmentName')}: "
          f"{sub_report.get('correctAnswers')}/{sub_report.get('totalQuestions')} correct, "
          f"{sub_report.get('obtainedMarks')}/{sub_report.get('totalMarks')} marks "
          f"({sub_report.get('percentage')}%)")
    print("\nDONE. Open the course in the app to see the video, PDF, quiz and practice problem.")
    print(f"  Course id : {cid}")
    print(f"  Student   : {student_email} / {student_pwd}")


if __name__ == "__main__":
    main()
