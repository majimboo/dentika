# Seeded Clinics and Users

## Clinics

### 1. Dentika (ID: 1) - Main Platform
**Address:** Unit 205 JY Square Mall, Lahug, Cebu City, 6000 Cebu
**Phone:** +63 32 520 8888
**Email:** info@dentika.com
**Website:** https://dentika.com
**Branches:** 1 (Main Branch)

**Users:**
- **Super Admin:** admin / admin (Admin User) - admin@dentika.com

**Patients:** No seeded patients (platform clinic for inventory management)

---

### 2. SmileCare Dental Clinic (ID: 2)
**Address:** G/F Ayala Center Cebu, Cebu Business Park, Cebu City, 6000 Cebu
**Phone:** +63 32 412 3456
**Email:** info@smilecare.ph
**Website:** https://smilecare.ph
**Branches:** 1 (Main Branch)

**Users:**
- **Admins:**
  - mcruz / admin (Maria Cruz) - maria.cruz@smilecare.ph
  - jreyes / admin (Jose Reyes) - jose.reyes@smilecare.ph
- **Doctors:**
  - asantos / admin (Ana Santos) - ana.santos@smilecare.ph
  - rgarcia / admin (Roberto Garcia) - roberto.garcia@smilecare.ph
- **Secretary:**
  - lvillanueva / admin (Luz Villanueva) - luz.villanueva@smilecare.ph
- **Assistants:**
  - cdela / admin (Carmen dela Cruz) - carmen.delacruz@smilecare.ph
  - jtanaka / admin (Juan Tanaka) - juan.tanaka@smilecare.ph

**Patients:** 20 seeded patients with Filipino names, addresses, and comprehensive medical information including:
- Personal details (name, email, phone, date of birth, gender, address)
- Emergency contacts with relationship details
- Medical information (blood type, allergies, medical conditions, current medications)
- Insurance details (PhilHealth, Maxicare, Intellicare)
- Notes for dental care considerations

---

### 3. Bright Smile Dental Center (ID: 3)
**Address:** 2/F Robinson's Galleria Cebu, Gen. Maxilom Avenue, Cebu City, 6000 Cebu
**Phone:** +63 32 234 5678
**Email:** info@brightsmile.ph
**Website:** https://brightsmile.ph
**Branches:** 2

#### Branch 1 - Main Branch (ID: 3)
**Address:** 2/F Robinson's Galleria Cebu, Gen. Maxilom Avenue, Cebu City, 6000 Cebu
**Phone:** +63 32 234 5678

#### Branch 2 - Lahug Branch (ID: 4)
**Address:** 3/F Ayala Terraces, Salinas Drive, Lahug, Cebu City, 6000 Cebu
**Phone:** +63 32 234 5679

**Users:**
- **Admins:**
  - mlopez / admin (Miguel Lopez) - miguel.lopez@brightsmile.ph
  - gjose / admin (Grace Jose) - grace.jose@brightsmile.ph
- **Doctors:**
  - fflores / admin (Ferdinand Flores) - ferdinand.flores@brightsmile.ph
  - maquino / admin (Michelle Aquino) - michelle.aquino@brightsmile.ph
- **Secretary:**
  - smendoza / admin (Sofia Mendoza) - sofia.mendoza@brightsmile.ph
- **Assistants:**
  - rramos / admin (Rosa Ramos) - rosa.ramos@brightsmile.ph
  - bcastro / admin (Benjamin Castro) - benjamin.castro@brightsmile.ph

**Patients:** 20 seeded patients with Filipino names, addresses, and comprehensive medical information including:
- Personal details (name, email, phone, date of birth, gender, address)
- Emergency contacts with relationship details
- Medical information (blood type, allergies, medical conditions, current medications)
- Insurance details (PhilHealth, Maxicare, Intellicare)
- Notes for dental care considerations

---

## Summary
- **Total Clinics:** 3
- **Total Branches:** 4
- **Total Users:** 15 (1 super admin + 14 clinic staff)
- **Total Patients:** 40 (20 for SmileCare + 20 for Bright Smile)

All passwords are set to "admin" for testing purposes.
All addresses are based in Cebu City, Philippines with local phone numbers.

## Database Schema Updates

### Consent Template Index Fix
- **Fixed:** `idx_consent_templates_code` to be unique per clinic
- **Changed:** From global unique `code` to composite unique `(code, clinic_id)`
- **Benefit:** Same consent template codes can now exist across different clinics
- **Index Name:** `idx_consent_templates_code_clinic`

This ensures that each clinic can have its own set of consent templates with the same codes (e.g., "EXTRACT-001") without conflicts.